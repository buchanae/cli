package inspect

import (
	"fmt"
  "log"
	"go/ast"
	"go/build"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"strings"
)

// TODO probably don't want Inspect writing to global log

func Inspect(packages []string) (*Package, error) {

	// Load the program.
	var conf loader.Config
	_, err := conf.FromArgs(packages, false)
	conf.ParserMode = parser.ParseComments
	// Try to be lenient about errors in the code.
	conf.TypeChecker.FakeImportC = true
	conf.TypeChecker.IgnoreFuncBodies = true
	conf.TypeChecker.DisableUnusedImportCheck = true
	conf.AllowErrors = true
  conf.TypeChecker.Error = func(e error) {}
	if err != nil {
    return nil, fmt.Errorf("configuring loader: %v", err)
	}

	prog, err := conf.Load()
	if err != nil {
    return nil, fmt.Errorf("loading program: %v", err)
	}

	initial := prog.InitialPackages()
	if len(initial) > 1 {
    return nil, fmt.Errorf("inspect doesn't understad multiple packages yet")
	}
	info := initial[0]

	p2, err := build.Default.Import(info.Pkg.Path(), ".", build.FindOnly)
	if err != nil {
    return nil, fmt.Errorf("finding import path: %v", err)
	}

	// Look for exported functions in the package.
	var funcs []*Func
	for _, file := range info.Files {

    filename := prog.Fset.Position(file.Package).Filename
    if !strings.HasSuffix(filename, "_cli.go") {
      continue
    }

		for _, dec := range file.Decls {
			if f, ok := dec.(*ast.FuncDecl); ok {
				if f.Name.IsExported() {
					funcs = append(funcs, &Func{
						Name:    f.Name.Name,
						Package: info.Pkg.Path(),
						Doc:     f.Doc.Text(),
					})
				}
			}
		}
	}

  // TODO inspect is reanalyzing the same option type many times,
  //      but it could probably cache the results on the first pass.
	// Gather information about the function arguments.
	scope := info.Pkg.Scope()
	for _, def := range funcs {
		obj := scope.Lookup(def.Name)
		z, ok := obj.(*types.Func)
    if !ok {
      // TODO this happens, but not exactly sure how yet.
      continue
    }
		sig := z.Type().(*types.Signature)

		params := sig.Params()
		for i := 0; i < params.Len(); i++ {
			p := params.At(i)

			isOpt := p.Name() == "opt"
			variadic := sig.Variadic() && i == params.Len()-1

			if isOpt {
				if variadic {
					continue
					// TODO error, opt cannot be variadic
				}

				// TODO what if "opt" is a string, int, etc?
				nt, ok := p.Type().(*types.Named)
				if !ok {
					continue
					// TODO error, opt must be a struct type
				}

				_, ok = nt.Underlying().(*types.Struct)
				if !ok {
					continue
					// TODO error, opt must be a struct type
				}

				tn := nt.Obj()

				defaultsName := "Default" + tn.Name()
				defObj := tn.Pkg().Scope().Lookup(defaultsName)
				// TODO validate the type of the defaults obj
				//      probably should allow only a function
				//      might consider a non-pointer value.
				def.HasDefaultOpts = defObj != nil

				def.Opts = walk(prog, nil, p.Type(), "")
				def.OptsType = nt

			} else {
				def.Args = append(def.Args, Arg{
					Name:     p.Name(),
					Type:     p.Type(),
					Variadic: variadic,
				})
			}
		}
	}

  return &Package{
    Name: info.Pkg.Name(),
    Path: info.Pkg.Path(),
    Dir: p2.Dir,
    Funcs: funcs,
  }, nil
}

type Package struct {
  Name string
  Path string
  Dir string
  Funcs []*Func
}

type Func struct {
	Name           string
	Package        string
	Doc            string
	Opts           []*Leaf
	OptsType       *types.Named
	HasDefaultOpts bool
	Args           []Arg
}

type Arg struct {
	Name     string
	Type     types.Type
	Variadic bool
}

// Leaf holds information about a leaf in a tree of struct fields.
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
type Leaf struct {
	// The path to the leaf field, e.g. "Root.Sub.SubOne"
	Key []string
	// The comment attached to the leaf, e.g. "Comment for SubOne field."
	Doc         string
	Type        types.Type
}

// walk recursively walks a struct, collecting leaf fields.
// See the `leaf` docs for more information about those fields.
func walk(prog *loader.Program, path []string, t types.Type, doc string) []*Leaf {
	var leaves []*Leaf

	switch t := t.(type) {
	case *types.Struct:

		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)
			if !f.Exported() {
				continue
			}

			subpath := path
			if !f.Anonymous() {
				subpath = newpathS(path, f.Name())
			}
			w := walk(prog, subpath, f.Type(), extractVarDoc(prog, f))
			leaves = append(leaves, w...)
		}

	case *types.Named:
    //return walk(prog, path, t.Underlying(), "")
		switch z := t.Underlying().(type) {
		case *types.Struct, *types.Named:
			return walk(prog, path, z, "")

    case *types.Pointer:

      switch el := z.Elem().(type) {
      case *types.Struct, *types.Named:
        return walk(prog, path, el, "")

      default:
        leaves = append(leaves, &Leaf{
          Key: path,
          Doc: doc,
          Type: t,
        })
      }

    case *types.Interface, *types.Basic, *types.Slice, *types.Map, *types.Array:
      leaves = append(leaves, &Leaf{
        Key: path,
        Doc: doc,
        Type: t,
      })

		default:
			// TODO what is going on here? not handling wrapper types?
			//      what is a value type?
      // TODO the path in this log message doesn't include the name of the root type.
      p := strings.Join(path, ".")
      log.Printf("skipping unhandled type at %q: %v\n", p, t)
			return nil
		}

  case *types.Pointer:
    return walk(prog, path, t.Elem(), doc)

	case *types.Basic, *types.Slice, *types.Map, *types.Array:
    leaves = append(leaves, &Leaf{
      Key:  path,
      Doc:  doc,
      Type: t,
    })

	default:
    p := strings.Join(path, ".")
	  log.Printf("skipping unhandled type at %q: %v\n", p, t)
	}
	return leaves
}

// extractVarDoc will attempt to return the code comment attached to a var,
// if it exists.
func extractVarDoc(prog *loader.Program, f types.Object) string {
	_, astpath, _ := prog.PathEnclosingInterval(f.Pos(), f.Pos())
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
