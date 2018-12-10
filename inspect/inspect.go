package inspect

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"os"
	"strings"
)

func Inspect(packages []string, verbose bool) (*Package, error) {

	// Load the program.
	var conf loader.Config
	_, err := conf.FromArgs(packages, false)
	conf.ParserMode = parser.ParseComments
	// Try to be lenient about errors in the code.
	conf.TypeChecker.FakeImportC = true
	conf.TypeChecker.IgnoreFuncBodies = true
	conf.TypeChecker.DisableUnusedImportCheck = true
	conf.AllowErrors = true
	if !verbose {
		conf.TypeChecker.Error = func(e error) {}
	}
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
		for _, dec := range file.Decls {
			if f, ok := dec.(*ast.FuncDecl); ok {
				if f.Name.IsExported() && strings.HasSuffix(f.Name.Name, "Cmd") {
					funcs = append(funcs, &Func{
						Name:    f.Name.Name,
						Package: info.Pkg.Path(),
						Doc:     f.Doc.Text(),
					})
				}
			}
		}
	}

	// Gather information about the function arguments.
	scope := info.Pkg.Scope()
	for _, def := range funcs {
		obj := scope.Lookup(def.Name)
		z := obj.(*types.Func)
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

				def.Opts = walkStruct(prog, verbose, nil, p.Type(), "")
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
	IsValueType bool
}

// walkStruct recursively walks a struct, collecting leaf fields.
// See the `leaf` docs for more information about those fields.
func walkStruct(prog *loader.Program, verbose bool, path []string, t types.Type, doc string) []*Leaf {
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
			w := walkStruct(prog, verbose, subpath, f.Type(), extractVarDoc(prog, f))
			leaves = append(leaves, w...)
		}

	case *types.Named:
		switch z := t.Underlying().(type) {
		case *types.Struct:
			return walkStruct(prog, verbose, path, z, "")
		default:
			// TODO hard-coded exception for funnel
			// TODO what is going on here? not handling wrapper types?
			//      what is a value type?
			fmt.Printf("UNHANDLED %#v\n", t)
			if t.String() == "github.com/ohsu-comp-bio/funnel/config.Duration" {
				leaves = append(leaves, &Leaf{
					Key:         path,
					Doc:         doc,
					Type:        t,
					IsValueType: true,
				})
				return leaves
			}

			if verbose {
				fmt.Fprintln(os.Stderr, "unhandled type", strings.Join(path, "."), t)
			}
			return nil
		}

	case *types.Basic:

		// Some basic types are not supported because they can't be defined as flags or config.
		switch t.Kind() {
		case types.Invalid, types.Uintptr, types.Complex64, types.Complex128, types.UnsafePointer,
			types.UntypedBool, types.UntypedInt, types.UntypedRune, types.UntypedFloat,
			types.UntypedComplex, types.UntypedString, types.UntypedNil:
			if verbose {
				fmt.Fprintln(os.Stderr, "unhandled type", strings.Join(path, "."), t)
			}
			return nil
		}

		leaves = append(leaves, &Leaf{
			Key:  path,
			Doc:  doc,
			Type: t,
		})

	case *types.Slice:
		if _, ok := t.Elem().(*types.Basic); ok {
			leaves = append(leaves, &Leaf{
				Key:  path,
				Doc:  doc,
				Type: t,
			})
		} else if verbose {
			fmt.Fprintln(os.Stderr, "unhandled type", strings.Join(path, "."), t)
		}

	default:
		if verbose {
			fmt.Fprintln(os.Stderr, "unhandled type", strings.Join(path, "."), t)
		}
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
