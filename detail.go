package cli

import (
	"bufio"
	"bytes"
	"go/doc"
	"strings"
)

// OptDetail holds information parsed from an option's doc string.
type OptDetail struct {
	Doc        string
	Synopsis   string
	Deprecated string
}

// ParseOptDetail parses an option's doc string for additional detail.
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

// CmdDetail holds information parsed from a command's doc string.
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

// ParseCmdDetail parses a command's doc string for additional information.
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

  // TODO should be building the new doc without annotations.
  var lines []string
	scan := bufio.NewScanner(bytes.NewBufferString(spec.Doc()))
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
    if line == det.Synopsis {
      continue
    }
		if strings.HasPrefix(line, "Name: ") {
			name := strings.TrimPrefix(line, "Name: ")
			parts := strings.Split(name, " ")
			setNamePath(parts)
      continue
		}
		if strings.HasPrefix(line, "Deprecated: ") {
			det.Deprecated = strings.TrimPrefix(line, "Deprecated: ")
      continue
		}
		if strings.HasPrefix(line, "Example: ") {
			det.Example = strings.TrimPrefix(line, "Example: ")
      continue
		}
		if line == "Hidden" {
			det.Hidden = true
      continue
		}
		if strings.HasPrefix(line, "Aliases: ") {
			det.Aliases = strings.Split(strings.TrimPrefix(line, "Aliases: "), " ")
      continue
		}
    lines = append(lines, line)
	}
  det.Doc = strings.TrimSpace(strings.Join(lines, "\n"))

	return det
}
