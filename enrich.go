package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"go/doc"
	"os"
	"reflect"
	"strings"
)

// Enrich fills in Cmd and Opt fields runtime
// by parsing names and doc strings, looking for annotations
// such as Synopsis, Deprecated, Example, etc.
func Enrich(cmd *Cmd) {
	enrichCmd(cmd)
	for _, opt := range cmd.Opts {
		enrichOpt(opt)
	}
}

// enrichOpt parses an option's doc string for additional detail.
func enrichOpt(opt *Opt) {
	opt.Synopsis = doc.Synopsis(opt.RawDoc)

	var lines []string
	scan := bufio.NewScanner(bytes.NewBufferString(opt.RawDoc))
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())

		// TODO "Name" is more tricky for options, because you might also want
		//      to override the name of intermediate parents, e.g. Foo.BarBAZ.Bat
		//      there's no Opt for BarBAZ to be parsed. Probably should rely
		//      on struct tags, on a runtime builder, and/or parse struct field
		//      docs in the code generator.
		//
		//      Also, struct tags can be reflected at runtime here in ParseOptDetail,
		//      if desired.
		//if strings.HasPrefix(line, "Name: ") {
		//det.Name = strings.TrimPrefix(line, "Name: ")
		//}

		switch {
		case line == opt.Synopsis:
		case line == "Hidden":
			opt.Hidden = true
		case strings.HasPrefix(line, "Deprecated: "):
			opt.Deprecated = strings.TrimPrefix(line, "Deprecated: ")
		default:
			lines = append(lines, line)
		}
	}
	opt.Doc = strings.TrimSpace(strings.Join(lines, "\n"))

	switch {
	case opt.DefaultValue == os.Stderr:
		opt.DefaultString = "os.Stderr"
	case opt.DefaultValue == os.Stdout:
		opt.DefaultString = "os.Stdout"
	case opt.DefaultValue == nil:
		opt.DefaultString = ""
		// TODO not super happy about including reflect just for this.
	case reflect.TypeOf(opt.DefaultValue).Kind() == reflect.Ptr:
		opt.DefaultString = ""
	default:
		opt.DefaultString = fmt.Sprintf("%v", opt.DefaultValue)
	}
}

// enrichCmd parses a command's doc string for additional information.
func enrichCmd(cmd *Cmd) {
	cmd.Synopsis = doc.Synopsis(cmd.RawDoc)

	setNamePath := func(parts []string) {
		if len(parts) == 0 {
			return
		}

		cmd.Path = nil
		for _, p := range parts {
			cmd.Path = append(cmd.Path, strings.ToLower(p))
		}
		cmd.Name = cmd.Path[len(cmd.Path)-1]
	}

	parts := splitIdent(cmd.RawName)
	setNamePath(parts)

	var lines []string
	scan := bufio.NewScanner(bytes.NewBufferString(cmd.RawDoc))
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())

		switch {
		case line == cmd.Synopsis:
		case strings.HasPrefix(line, "Name: "):
			name := strings.TrimPrefix(line, "Name: ")
			parts := strings.Split(name, " ")
			setNamePath(parts)
		case strings.HasPrefix(line, "Deprecated: "):
			cmd.Deprecated = strings.TrimPrefix(line, "Deprecated: ")
			// TODO example could be multi line
		case strings.HasPrefix(line, "Example: "):
			cmd.Example = strings.TrimPrefix(line, "Example: ")
		case line == "Hidden":
			cmd.Hidden = true
		case strings.HasPrefix(line, "Aliases: "):
			cmd.Aliases = strings.Split(strings.TrimPrefix(line, "Aliases: "), " ")
		default:
			lines = append(lines, line)
		}
	}
	cmd.Doc = strings.TrimSpace(strings.Join(lines, "\n"))
}
