/*
roger generates Go code to be used by the github.com/buchanae/roger/roger library.
*/
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
  var root string
  alias := mapVar{}
  ignore := sliceVar{}

  flag.StringVar(&outpath, "out", outpath, "File to write generated output to. Required.")
  flag.StringVar(&root, "root", root, "Name of the struct to inspect. Required.")
  flag.Var(&ignore, "i", "Ignore these fields.")
  flag.Var(alias, "a", `Alias these fields, e.g. "short=Path.To.Struct.Field".`)
  flag.BoolVar(&verbose, "v", verbose, "Verbose logging.")
  flag.Parse()

  if outpath == "" {
    fmt.Fprintln(os.Stderr, "-output is required")
    fmt.Fprintln(os.Stderr, "usage: roger -root Config -output out.go ./inputs [...]")
    flag.PrintDefaults()
    os.Exit(1)
  }

  if root == "" {
    fmt.Fprintln(os.Stderr, "-root is required")
    fmt.Fprintln(os.Stderr, "usage: roger -root Config -output out.go ./inputs [...]")
    flag.PrintDefaults()
    os.Exit(1)
  }

  // Load the program.
  var conf loader.Config

  _, err := conf.FromArgs(flag.Args(), false)
  conf.ParserMode = parser.ParseComments

  // Try to be lenient about errors in the code.
  conf.TypeChecker.FakeImportC = true
  conf.TypeChecker.IgnoreFuncBodies = true
  conf.TypeChecker.DisableUnusedImportCheck = true
  conf.AllowErrors = true

	if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
	}

  prog, err := conf.Load()
	if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
	}

  // Find the root object config structure to inspect.
  var rootobj types.Object
  var name string

  for _, info := range prog.InitialPackages() {
    rootobj = info.Pkg.Scope().Lookup(root)
    if rootobj != nil {
      name = info.Pkg.Name()
      break
    }
  }

  if rootobj == nil {
    fmt.Fprintf(os.Stderr, "Can't find root struct named '%s'\n", root)
    os.Exit(1)
  }

  // Walk the config structure, building a list of key/value items.
  leaves := walkStruct(prog, nil, rootobj.Type(), verbose)

  // Filter leaves based on "-ignore" command line flag.
  var filtered []*leaf
  for i, n := range leaves {

    shouldIgnore := false
    k := strings.Join(n.Key, ".")
    for _, ig := range ignore {
      if strings.HasPrefix(k, ig) {
        shouldIgnore = true
        break
      }
    }

    if !shouldIgnore {
      filtered = append(filtered, leaves[i])
    }
  }
  leaves = filtered

  // Generate the code.
  var b bytes.Buffer
  tpl.Execute(&b, map[string]interface{}{
    "Pkgname": name,
    "Leaves": leaves,
    "Alias": alias,
  })
  s, err := format.Source(b.Bytes())
  if err != nil {
    panic(err)
  }

  // Write the code.
  out, err := os.Create(outpath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: %s", err)
    os.Exit(1)
  }
  defer out.Close()
  fmt.Fprintln(out, string(s))
}

// tpl is used to render the generated code. See tpl.go
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

// leaf holds information about a leaf in a tree of struct fields.
// For example:
// 
//   type Root struct {
//     RootOne string
//     Sub struct {
//       // Comment for SubOne field.
//       SubOne string
//     }
//   }
//
// Root.RootOne and Root.Sub.SubOne are leaves.
type leaf struct {
  // The path to the leaf field, e.g. "Root.Sub.SubOne"
  Key []string
  // The comment attached to the leaf, e.g. "Comment for SubOne field."
  Doc string
}

// walkStruct recursively walks a struct, collecting leaf fields.
// See the `leaf` docs for more information about those fields.
func walkStruct(prog *loader.Program, path []string, t types.Type, verbose bool) []*leaf {
  switch t := t.(type) {

  case *types.Struct:
    var leaves []*leaf

    for i := 0; i < t.NumFields(); i++ {
      f := t.Field(i)
      if !f.Exported() {
        continue
      }

      subpath := newpathS(path, f.Name())

      if w := walkStruct(prog, subpath, f.Type(), verbose); w != nil {
        leaves = append(leaves, w...)
      } else {
        leaves = append(leaves, &leaf{
          Key: subpath,
          Doc: extractVarDoc(prog, f),
        })
      }
    }
    return leaves

  case *types.Named:
    return walkStruct(prog, path, t.Underlying(), verbose)

  case *types.Basic:
  default:
    if verbose {
      fmt.Fprintln(os.Stderr, "unhandled type", strings.Join(path, "."), t)
    }
  }
  return nil
}

// extractVarDoc will attempt to return the code comment attached to a var,
// if it exists.
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

// extractFieldDoc will attempt to return the code comment attached to an ast.Node,
// if it exists.
func extractFieldDoc(n ast.Node) string {
  switch n := n.(type) {
  case *ast.Field:
    if n.Doc != nil {
      return n.Doc.Text()
    }
  }
  return ""
}

// newpathS helps copy a slice of strings representing the path to a struct field.
func newpathS(base []string, add ...string) []string {
  path := append([]string{}, base...)
  return append(path, add...)
}

// mapVar is used to capture the `alias` command line flag.
type mapVar map[string]string
func (m mapVar) String() string {
  return fmt.Sprintf("%#v", m)
}
func (m mapVar) Get() interface{} {
  return m
}
func (m mapVar) Set(s string) error {
  sp := strings.Split(s, "=")
  if len(sp) == 2 {
    m[sp[0]] = sp[1]
    return nil
  }
  return fmt.Errorf("unrecognized alias: %s", s)
}

// sliceVar is used to capture the `ignore` command line flag.
type sliceVar []string
func (sv *sliceVar) String() string {
  return fmt.Sprintf("%#v", *sv)
}
func (sv *sliceVar) Get() interface{} {
  return *sv
}
func (sv *sliceVar) Set(s string) error {
  *sv = append(*sv, s)
  return nil
}
