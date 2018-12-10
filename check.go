package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"go/doc"
	"os"
	"strings"
)

type ErrFatal struct {
	err error
}

func (e ErrFatal) Error() string {
	return e.err.Error()
}

func Fatal(msg string, args ...interface{}) {
	panic(ErrFatal{fmt.Errorf(msg, args...)})
}

func Check(err error) {
	if err != nil {
		panic(ErrFatal{err})
	}
}

func Run(spec CmdSpec, args []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if x, ok := r.(error); ok {
				err = x
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	spec.Run(args)
	return
}

func Open(path string) *os.File {
	fh, err := os.Open(path)
	Check(err)
	return fh
}

func CheckArgs(args []string, specs []ArgSpec) {
	if len(specs) == 0 && len(args) > 0 {
		Fatal("unexpected args %v", args)
	}

	variadic := specs[len(specs)-1].Variadic
	if variadic {
		min := len(specs) - 1
		if len(args) < min {
			Fatal("expected at least %d arg", min)
		}
	} else {
		if len(args) != len(specs) {
			Fatal("expected exactly %d args", len(specs))
		}
	}
}

func CoerceString(arg string) string {
	return arg
}

func CoerceInts(args []string) []int {
	return nil
}

func CoerceInt(args string) int {
	return 0
}

type OptDetail struct {
	Doc        string
	Synopsis   string
	Deprecated string
}

func ParseOptDetail(spec OptSpec) OptDetail {
	det := OptDetail{
		Doc:      spec.Doc,
		Synopsis: doc.Synopsis(spec.Doc),
	}

	scan := bufio.NewScanner(bytes.NewBufferString(spec.Doc))
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		// TODO this is more tricky for options, because you might also want
		//      to override the name of intermediate parents, e.g. Foo.BarBAZ.Bat
		//      there's no OptSpec for BarBAZ to be parsed. Probably should rely
		//      on struct tags, on a runtime builder, and/or parse struct field
		//      docs in the code generator.
		//
		//      Also, struct tags can be reflected at runtime here in ParseOptDetail,
		//      if desired.
		//if strings.HasPrefix(line, "Name: ") {
		//det.Name = strings.TrimPrefix(line, "Name: ")
		//}
		if strings.HasPrefix(line, "Deprecated: ") {
			det.Deprecated = strings.TrimPrefix(line, "Deprecated: ")
		}
	}
	return det
}

type CmdDetail struct {
	Name       string
	Path       []string
	Doc        string
	Synopsis   string
	Example    string
	Deprecated string
	Hidden     bool
	Aliases    []string
}

func ParseCmdDetail(spec CmdSpec) CmdDetail {

	det := CmdDetail{
		Doc:      spec.Doc(),
		Synopsis: doc.Synopsis(spec.Doc()),
	}

	setNamePath := func(parts []string) {
		if len(parts) == 0 {
			return
		}

		for _, p := range parts {
			det.Path = append(det.Path, strings.ToLower(p))
		}
		det.Name = det.Path[len(det.Path)-1]
	}

	parts := SplitIdent(spec.Name())
	setNamePath(parts)

	scan := bufio.NewScanner(bytes.NewBufferString(spec.Doc()))
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if strings.HasPrefix(line, "Name: ") {
			name := strings.TrimPrefix(line, "Name: ")
			parts := strings.Split(name, " ")
			setNamePath(parts)
		}
		if strings.HasPrefix(line, "Deprecated: ") {
			det.Deprecated = strings.TrimPrefix(line, "Deprecated: ")
		}
		if strings.HasPrefix(line, "Example: ") {
			det.Example = strings.TrimPrefix(line, "Example: ")
		}
		if line == "Hidden" {
			det.Hidden = true
		}
		if strings.HasPrefix(line, "Aliases: ") {
			det.Aliases = strings.Split(strings.TrimPrefix(line, "Aliases: "), " ")
		}
	}

	return det
}
