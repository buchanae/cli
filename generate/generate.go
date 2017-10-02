package main

import (
  "text/template"
  "go/ast"
  "go/parser"
  "go/types"
  "fmt"
  "os"
  "golang.org/x/tools/go/loader"
  "strings"
)

func main() {
  ParseComments()
}

var tpl = template.Must(template.New("gen").
  Funcs(map[string]interface{}{
    "join": strings.Join,
  }).
  Parse(`
package main

func (c *Config) Set(k string, v interface{}) error {
  var ptrs = map[string]interface{}{
    {{ range .Nodes -}}
      "{{ join .Key "." }}": &c.{{ join .Key "." }},
    {{ end }}
  }
  ptrs[k] = v
  return nil
}
  `),
)

func ParseComments() {

  var conf loader.Config
  _, err := conf.FromArgs([]string{"github.com/buchanae/roger"}, false)
  conf.ParserMode = parser.ParseComments

	if err != nil {
		panic(err)
	}

  prog, err := conf.Load()
	if err != nil {
		panic(err)
	}

  pkg := prog.Package("github.com/buchanae/roger")
  co := pkg.Pkg.Scope().Lookup("Config")
  nodes := walkDocs(prog, nil, co.Type())

  tpl.Execute(os.Stdout, map[string]interface{}{
    "Nodes": nodes,
  })

  /*
  for _, n := range nodes {
    fmt.Printf("fs.c.%s\n%s\n\n", strings.Join(n.Key, "."), doc.Synopsis(n.Doc))
  }
  */
}

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

func walkDocs(prog *loader.Program, path []string, t types.Type) []*docnode {
  switch t := t.(type) {

  case *types.Struct:
    var nodes []*docnode

    for i := 0; i < t.NumFields(); i++ {
      f := t.Field(i)
      if !f.Exported() {
        continue
      }

      subpath := newpathS(path, f.Name())

      if w := walkDocs(prog, subpath, f.Type()); w != nil {
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
    return walkDocs(prog, path, t.Underlying())

  case *types.Basic:
  default:
    fmt.Fprintln(os.Stderr, "unknown type", t)
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
