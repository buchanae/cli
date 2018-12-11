package inspect

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
  "go/doc"
  "log"
	"os"
	"path/filepath"
	"text/template"
	"strings"
)

type ErrGofmt error

func Generate(pkg *Package, tpl *template.Template) (err error) {

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, TemplateVars(pkg))
	if err != nil {
    err = fmt.Errorf("executing template: %v", err)
    return
	}

	// Try to format (tidy) the source code, but if it fails just skip it.
	by := buf.Bytes()
	s, err := format.Source(by)
	if err != nil {
    err = ErrGofmt(fmt.Errorf("Go code formatting failed, producing unformatted code. %v", err))
	} else {
		by = s
	}

	// Write the code.
	outPath := filepath.Join(pkg.Dir, "generated_specs.go")
	out, err := os.Create(outPath)
	if err != nil {
    err = fmt.Errorf("creating output file %q: %v", outPath, err)
    return
	}
	defer out.Close()

  log.Printf("generated file %s\n", outPath)

	fmt.Fprintln(out, string(by))
  return
}

type uniqImports map[string]string

func (u uniqImports) Uniq(pkgname string) string {
  try := pkgname
  i := 1
  for {
    _, ok := u[pkgname]
    if !ok {
      break
    }
    try = fmt.Sprint("%s%d", pkgname, i)
  }
  return try
}

func TemplateVars(pkg *Package) map[string]interface{} {
	var defs []tplVars

	imports := uniqImports{
		"cli": "github.com/buchanae/cli",
	}

	for _, def := range pkg.Funcs {
    name := def.Name
		vars := tplVars{
			FuncName:     name,
			FuncNamePriv: makePrivate(name),
			Doc:          def.Doc,
		}

		for i, arg := range def.Args {
      typeName := arg.Type.String()
      if nt, ok := arg.Type.(*types.Named); ok {
        tn := nt.Obj()
        path := tn.Pkg().Path()
        name := tn.Name()
        if path != def.Package {
          pkgname := imports.Uniq(tn.Pkg().Name())
          imports[pkgname] = path
          typeName = pkgname + "." + name
        }
      }

			vars.Args = append(vars.Args, argVars{
				Idx:        i,
				Name:       arg.Name,
				Type:       typeName,
				Variadic:   arg.Variadic,
			})
		}
		vars.HasArgs = len(vars.Args) > 0

		if def.Opts != nil {
			vars.HasOpts = true
			vars.HasDefaultOpts = def.HasDefaultOpts

			tn := def.OptsType.Obj()
			path := tn.Pkg().Path()
			name := tn.Name()

			if path == def.Package {
				vars.OptsType = name
				vars.DefaultOptsName = "Default" + name + "()"
			} else {
				pkgname := imports.Uniq(tn.Pkg().Name())
				imports[pkgname] = path
				vars.OptsType = pkgname + "." + name
				vars.DefaultOptsName = pkgname + ".Default" + name + "()"
			}
		}

		for _, opt := range def.Opts {
			vars.Opts = append(vars.Opts, optVars{
				Key:        opt.Key,
				KeyJoined:  strings.Join(opt.Key, "."),
				Type:       opt.Type.String(),
				Doc:        opt.Doc,
				Synopsis:   doc.Synopsis(opt.Doc),
				Deprecated: "",
				Hidden:     false,
			})
		}
		defs = append(defs, vars)
	}

	return map[string]interface{}{
		"Funcs":   defs,
		"Package": pkg.Name,
		"Imports": imports,
	}
}

func makePrivate(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

type optVars struct {
	Key                       []string
	Doc, Synopsis, Deprecated string
	KeyJoined                 string
	Hidden                    bool
	Type                      string
}

type tplVars struct {
	FuncName, FuncNamePriv, Synopsis, Doc, Example, Deprecated string
	Aliases                                                    []string
	Hidden                                                     bool

	HasOpts         bool
	HasDefaultOpts  bool
	DefaultOptsName string
	OptsType        string
	Opts            []optVars

	HasArgs bool
	Args    []argVars
}

type argVars struct {
	Idx        int
	Name       string
	Type       string
	Variadic   bool
}
