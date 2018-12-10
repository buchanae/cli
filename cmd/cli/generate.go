package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/format"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	//"github.com/sanity-io/litter"
	"strings"
)

func init() {
	log.SetFlags(0)
}

var verbose bool

type Arg struct {
	Name     string
	Type     types.Type
	Variadic bool
}

type FuncDef struct {
	Name           string
	Package        string
	Doc            string
	Opts           []*Leaf
	OptsType       *types.Named
	HasDefaultOpts bool
	Args           []Arg
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

func main() {
	var err error
	verbose = false

	flag.BoolVar(&verbose, "v", verbose, "Verbose logging.")
	flag.Parse()

	tplBytes, err := ioutil.ReadFile("_tpl.go")
	if err != nil {
		log.Fatal(err)
	}
	tplstr := string(tplBytes)

	tpl, err := template.New("gen").Funcs(tplfuncs).Parse(tplstr)
	if err != nil {
		log.Fatal(err)
	}

	// Load the program.
	var conf loader.Config

	_, err = conf.FromArgs(flag.Args(), false)
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
		log.Fatal(err)
	}

	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	var funcs []*FuncDef

	initial := prog.InitialPackages()
	if len(initial) > 1 {
		log.Fatal("cli generator doesn't understand multiple packages (yet?)")
	}
	info := initial[0]

	p2, err := build.Default.Import(info.Pkg.Path(), ".", build.FindOnly)
	if err != nil {
		log.Fatal(err)
	}
	rootDir := p2.Dir

	// Look for exported functions in the package.
	for _, file := range info.Files {
		for _, dec := range file.Decls {
			if f, ok := dec.(*ast.FuncDecl); ok {
				if f.Name.IsExported() && strings.HasSuffix(f.Name.Name, "Cmd") {
					funcs = append(funcs, &FuncDef{
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

				def.Opts = walkStruct(prog, nil, p.Type(), "")
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

	//fmt.Fprintln(os.Stderr, litter.Sdump(funcs))

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, prepareTplVars(info.Pkg.Name(), funcs))
	if err != nil {
		log.Println(err)
	}

	// Try to format (tidy) the source code, but if it fails just skip it.
	by := buf.Bytes()
	s, err := format.Source(by)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Go code formatting failed, producing unformatted code. %s\n", err)
	} else {
		by = s
	}

	// Write the code.
	outPath := filepath.Join(rootDir, "generated_cli.go")
	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	fmt.Fprintln(out, string(by))
}

func makePrivate(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func prepareTplVars(pkgName string, funcs []*FuncDef) map[string]interface{} {
	var defs []tplVars

	imports := map[string]string{
		"cli": "github.com/buchanae/cli",
	}

	for _, def := range funcs {
		name := strings.TrimSuffix(def.Name, "Cmd")
		vars := tplVars{
			FuncName:     name,
			FuncNamePriv: makePrivate(name),
			Doc:          def.Doc,
		}

		for i, arg := range def.Args {

			coerceType := ""
			switch arg.Type.String() {
			case "string":
				coerceType = "String"
			case "[]string":
				coerceType = "Strings"
			case "int":
				coerceType = "Int"
			case "[]int":
				coerceType = "Ints"
			}

			vars.Args = append(vars.Args, argVars{
				Idx:        i,
				Name:       arg.Name,
				Type:       arg.Type.String(),
				CoerceType: coerceType,
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

				pkgname := tn.Pkg().Name()
				i := 1
				for {
					_, ok := imports[pkgname]
					if !ok {
						break
					}
					pkgname = fmt.Sprint("%s%d", tn.Pkg().Name(), i)
				}

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
		"Package": pkgName,
		"Imports": imports,
	}
}

/*
type CmdSpec struct {
  Name, Synopsis, Doc, Example, Deprecated string
  Aliases []string
  Hidden bool
}
*/

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
	CoerceType string
}

// Template functions
var tplfuncs = map[string]interface{}{
	// Return the pflag flag type function prefix,
	// e.g. "String" for a "StringVar" flag.
	// https://godoc.org/github.com/spf13/pflag
	"pflagType": func(l *Leaf) string {
		switch l.Type.String() {
		case "string", "int", "int64", "int32", "int16", "int8", "bool", "float32", "float64",
			"uint", "uint16", "uint32", "uint64", "uint8":
			return strings.Title(l.Type.String())
		case "[]string":
			return "StringSlice"
		}
		return "Unknown"
	},
}

// walkStruct recursively walks a struct, collecting leaf fields.
// See the `leaf` docs for more information about those fields.
func walkStruct(prog *loader.Program, path []string, t types.Type, doc string) []*Leaf {
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
			w := walkStruct(prog, subpath, f.Type(), extractVarDoc(prog, f))
			leaves = append(leaves, w...)
		}

	case *types.Named:
		switch z := t.Underlying().(type) {
		case *types.Struct:
			return walkStruct(prog, path, z, "")
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
