package main

import (
  "bytes"
  "flag"
  "text/template"
  "go/ast"
  "go/doc"
  "go/parser"
  "go/types"
  "go/format"
  "fmt"
  "os"
  "golang.org/x/tools/go/loader"
  "strings"
)

func main() {
  var outpath string
  var verbose bool
  flag.StringVar(&outpath, "out", outpath, "File to write generated output to.")
  flag.BoolVar(&verbose, "v", verbose, "Verbose logging.")
  flag.Parse()

  if outpath == "" {
    fmt.Fprintln(os.Stderr, "usage: roger -output out.go ./inputs ...")
    flag.PrintDefaults()
    os.Exit(1)
  }

  var conf loader.Config

  _, err := conf.FromArgs(flag.Args(), false)
  conf.ParserMode = parser.ParseComments

  // Try to be lenient about errors in the code.
  conf.TypeChecker.FakeImportC = true
  conf.TypeChecker.IgnoreFuncBodies = true
  conf.TypeChecker.DisableUnusedImportCheck = true
  conf.AllowErrors = true

	if err != nil {
		panic(err)
	}

  prog, err := conf.Load()
	if err != nil {
		panic(err)
	}

  // Find the target config structure to inspect.
  var target types.Object
  var name string

  for _, info := range prog.InitialPackages() {
    target = info.Pkg.Scope().Lookup("Config")
    if target != nil {
      name = info.Pkg.Name()
      break
    }
  }

  if target == nil {
    panic("can't find Config")
  }

  // Walk the config structure, building a list of key/value items.
  nodes := walkDocs(prog, nil, target.Type(), verbose)

  // Generate the configuration code to stdout.
  var b bytes.Buffer
  tpl.Execute(&b, map[string]interface{}{
    "Pkgname": name,
    "Nodes": nodes,
  })
  s, err := format.Source(b.Bytes())
  if err != nil {
    panic(err)
  }

  out, err := os.Create(outpath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: %s", err)
    os.Exit(1)
  }
  defer out.Close()
  fmt.Fprintln(out, string(s))

  /*
  for _, n := range nodes {
    fmt.Printf("fs.c.%s\n%s\n\n", strings.Join(n.Key, "."), doc.Synopsis(n.Doc))
  }
  */
}

var tpl = template.Must(template.New("gen").
  Funcs(map[string]interface{}{
    "join": func(s []string) string {
      return strings.Join(s, ".")
    },
    "synopsis": func(s string) string {
      return doc.Synopsis(s)
    },
    "keysyn": func(s []string) string {
      return fmt.Sprintf("%#v", s)
    },
  }).
  Parse(rawtpl),
)

type docnode struct {
  Key []string
  Doc string
}

func extractVarDoc(prog *loader.Program, f *types.Var) string {
  _, astpath, _ := prog.PathEnclosingInterval(f.Pos(), f.Pos())
  // TODO something here is wrong. This will search all the way up the path.
  for _, n := range astpath {
    d := extractFieldDoc(n)
    if d != "" {
      return d
    }
  }
  return ""
}

func walkDocs(prog *loader.Program, path []string, t types.Type, verbose bool) []*docnode {
  switch t := t.(type) {

  case *types.Struct:
    var nodes []*docnode

    for i := 0; i < t.NumFields(); i++ {
      f := t.Field(i)
      if !f.Exported() {
        continue
      }

      subpath := newpathS(path, f.Name())

      if w := walkDocs(prog, subpath, f.Type(), verbose); w != nil {
        nodes = append(nodes, w...)
      } else {
        nodes = append(nodes, &docnode{
          Key: subpath,
          Doc: extractVarDoc(prog, f),
        })
      }
    }
    return nodes

  case *types.Named:
    return walkDocs(prog, path, t.Underlying(), verbose)

  case *types.Basic:
  default:
    if verbose {
      fmt.Fprintln(os.Stderr, "unhandled type", t)
    }
  }
  return nil
}

func extractFieldDoc(n ast.Node) string {
  switch n := n.(type) {
  case *ast.Field:
    if n.Doc != nil {
      return n.Doc.Text()
    }
  }
  return ""
}

func newpathS(base []string, add ...string) []string {
  path := append([]string{}, base...)
  return append(path, add...)
}
