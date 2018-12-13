package cli

import (
	"bufio"
	"bytes"
	"go/doc"
	"strings"
)

type Spec interface {
  Cmd() *Cmd
  Run()
}

type Cmd struct {
  FuncName string
  RawDoc string
  Args []*Arg
  Opts []*Opt

	Name       string
	Path       []string
	Doc        string
	Synopsis   string
	Example    string
	Deprecated string
	Hidden     bool
	Aliases    []string
}

// Opt defines an option of a Cmd.
type Opt struct {
	Key   []string
	RawDoc   string
	Doc        string
  Short      string
	Synopsis   string
  // TODO parse hidden from doc
  Hidden     bool
	Deprecated string
  Type string
	Value interface{}
  DefaultValue interface{}
  IsSet bool
}

// Arg defines a positional argument of a Cmd.
type Arg struct {
	Name     string
	Type     string
	Variadic bool
  Value interface{}
}

// Enrich parses Cmd.RawDoc and extracts fields
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

		if strings.HasPrefix(line, "Deprecated: ") {
			opt.Deprecated = strings.TrimPrefix(line, "Deprecated: ")
		}
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

	parts := SplitIdent(cmd.FuncName)
	setNamePath(parts)

  var lines []string
	scan := bufio.NewScanner(bytes.NewBufferString(cmd.RawDoc))
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
    if line == cmd.Synopsis {
      continue
    }
		if strings.HasPrefix(line, "Name: ") {
			name := strings.TrimPrefix(line, "Name: ")
			parts := strings.Split(name, " ")
			setNamePath(parts)
      continue
		}
		if strings.HasPrefix(line, "Deprecated: ") {
			cmd.Deprecated = strings.TrimPrefix(line, "Deprecated: ")
      continue
		}
		if strings.HasPrefix(line, "Example: ") {
			cmd.Example = strings.TrimPrefix(line, "Example: ")
      continue
		}
		if line == "Hidden" {
			cmd.Hidden = true
      continue
		}
		if strings.HasPrefix(line, "Aliases: ") {
			cmd.Aliases = strings.Split(strings.TrimPrefix(line, "Aliases: "), " ")
      continue
		}
    lines = append(lines, line)
	}
  cmd.Doc = strings.TrimSpace(strings.Join(lines, "\n"))
}
