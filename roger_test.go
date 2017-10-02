package roger

import (
  "fmt"
  "flag"
  "go/ast"
  "go/parser"
  "go/types"
  //"go/doc"
  //"go/token"
  "testing"
  "os"
  "io/ioutil"
  "encoding/json"
  "github.com/kr/pretty"
  "golang.org/x/tools/go/loader"
	"github.com/ghodss/yaml"
  "github.com/buchanae/roger/example/server"
  "github.com/buchanae/roger/example/worker"
  "github.com/buchanae/roger/example/scheduler"
  "github.com/buchanae/roger/example/logger"
  "github.com/buchanae/roger/example/dynamo"
  "github.com/spf13/cast"
  "github.com/alecthomas/units"
  "time"
  "reflect"
  "strings"
)

/* Why

- Outside of the cli code, codebase deals directly is stucts. i.e. config/arg code
  does not leak into the codebase.

- Defaults are also defined by structs (i.e. pulled from the codebase) instead of
  being cli/arg specific

- Help/docs come from the code

- cli/arg/config code should live *around* your code, not *in* it.
*/

/*
Done:
- env prefix
- flag/env var naming style ServiceName -> service-name or service_name or service.name
- set defaults from a struct
- auto inspect
- flag generation
- read from flag, env, yaml, json
- support time.Duration in yaml, json, env
- report unknown fields
- alias/link/source field value from another field
- ignore/hide fields
- define short fields
- validation interface
- choices + validation
- support byte size
- support SI prefix (K, G, M, etc)

TODO:
- help/docs from comments
- printing config, but only non-defaults
- manage editing config file
- pluggable sources
- sets of default configurations
- slice of choices
- improve stringSlice.String() format
- dump json
- handle map[string]string via "key=value" flag value
- explore "storage.local.allowed_dirs.append"
- pull fieldname from json tag
- recognize misspelled env var
- case sensitivity

Complex:
- reloading
- multiple config files with merging
*/

/* Questions
- how to handle pointers? cycles?
- how are slices handled in env vars?
- how are slices of structs handled in flags?
- how to handle unknown type wrappers, e.g. type Foo int64
- want to hide all fields below a given prefix?
- how to handle string slice from env? comma sep? make it consistent with flag?
*/

func DefaultConfig() Config {
  return Config{
    Server: server.DefaultConfig(),
    Worker: worker.DefaultConfig(),
    Scheduler: scheduler.DefaultConfig(),
    Log: logger.DefaultConfig(),
    Dynamo: dynamo.DefaultConfig(),
  }
}

type Config struct {
  Server server.Config
  Worker worker.Config
  Scheduler scheduler.Config
  Log logger.Config
  Dynamo dynamo.Config
}

type Validator interface {
  Validate() []error
}

type leaf struct {
  Path []string
  Type reflect.Type
  Value reflect.Value
  Addr interface{}
}

func (l *leaf) String() string {
  return l.Value.String()
}

func (l *leaf) Set(s string) error {
  return l.Coerce(s)
}

func (l *leaf) Get() interface{} {
  return l.Value.Interface()
}

func (l *leaf) Coerce(val interface{}) error {
  var casted interface{}
  var err error

  switch l.Value.Interface().(type) {
  case int:
    casted, err = cast.ToIntE(val)
  case int64:
    casted, err = cast.ToInt64E(val)
  case int32:
    casted, err = cast.ToInt32E(val)
  case float32:
    casted, err = cast.ToFloat32E(val)
  case float64:
    casted, err = cast.ToFloat64E(val)
  case bool:
    casted, err = cast.ToBoolE(val)
  case string:
    casted, err = cast.ToStringE(val)
  case []string:
    casted, err = cast.ToStringSliceE(val)
  case units.MetricBytes:
    if s, ok := val.(string); ok {
      casted, err = units.ParseMetricBytes(s)
    }
  case time.Duration:
    casted, err = cast.ToDurationE(val)
  default:
    return fmt.Errorf("unknown source value", l.Path, val)
  }

  if err != nil {
    return fmt.Errorf("error casting", l.Path, val, err)
  }

  l.Value.Set(reflect.ValueOf(casted))
  return nil
}



func Inspect(i interface{}, hide []string) *tree {
  t := reflect.TypeOf(i)
  v := reflect.ValueOf(i)

  if v.Kind() != reflect.Ptr || v.IsNil() {
    panic("must be non-nil pointer type")
  }

  // TODO check that it's a struct type

  tr := tree{
    leaves: map[string]*leaf{},
    st: t.Elem(),
    sv: v.Elem(),
    hide: hide,
    namer: flagname,
  }
  tr.inspect(nil)
  return &tr
}

type tree struct {
  leaves map[string]*leaf
  hide []string
  st reflect.Type
  sv reflect.Value
  ignoreEmpty bool
  namer func(path []string) string
}

func (tr *tree) pathname(path []int) []string {
  var name []string
  for i := 0; i < len(path); i++ {
    name = append(name, tr.st.FieldByIndex(path[:i+1]).Name)
  }
  return name
}

func (tr *tree) validate(base []int) (errs []error) {
  t := tr.sv.FieldByIndex(base)

  for j := 0; j < t.NumField(); j++ {
    path := newpathI(base, j)
    if tr.shouldhide(path) {
      continue
    }

    fv := tr.sv.FieldByIndex(path)
    name := tr.namer(tr.pathname(path))

    if x, ok := fv.Interface().(Validator); ok {
      for _, err := range x.Validate() {
        errs = append(errs, fmt.Errorf("%s: %s", name, err))
      }
    }

    if fv.Kind() == reflect.Struct {
      errs = append(errs, tr.validate(path)...)
    }
  }
  return
}

func (tr *tree) dump(base []int) {
  t := tr.sv.FieldByIndex(base)

  for j := 0; j < t.NumField(); j++ {
    indent := strings.Repeat("  ", len(base))
    path := newpathI(base, j)
    if tr.shouldhide(path) {
      continue
    }

    ft := tr.st.FieldByIndex(path)
    fv := tr.sv.FieldByIndex(path)

    // Ignore zero values if ignoreEmpty is true.
    zero := reflect.Zero(ft.Type)
    eq := reflect.DeepEqual(zero.Interface(), fv.Interface())
    if tr.ignoreEmpty && eq {
      continue
    }

    switch fv.Kind() {
    case reflect.Struct:
      fmt.Printf("%s%s:\n", indent, ft.Name)
      tr.dump(path)

    default:
      fmt.Printf("%s%s: %v\n", indent, ft.Name, fv)
    }
  }
}

func (tr *tree) shouldhide(path []int) bool {
  pathname := tr.namer(tr.pathname(path))
  for _, h := range tr.hide {
    if strings.HasPrefix(pathname, h) {
      return true
    }
  }
  return false
}

func (tr *tree) inspect(base []int) {
  t := tr.sv.FieldByIndex(base)

  for j := 0; j < t.NumField(); j++ {
    path := newpathI(base, j)
    if tr.shouldhide(path) {
      continue
    }
    fv := tr.sv.FieldByIndex(path)

    switch fv.Kind() {
    case reflect.Struct:
      tr.inspect(path)

    default:
      name := tr.namer(tr.pathname(path))
      tr.leaves[name] = &leaf{
        Path: tr.pathname(path),
        Type: fv.Type(),
        Value: fv,
        Addr: fv.Addr().Interface(),
      }
    }
  }
}

func (tr *tree) LoadEnv(envname func(path []string) string) {
  for _, l := range tr.leaves {
    v := os.Getenv(envname(l.Path))
    if v != "" {
      l.Coerce(v)
    }
  }
}

type collector struct {}
func (c collector) Visit(n ast.Node) ast.Visitor {
  fmt.Println("collect", reflect.TypeOf(n))
  pretty.Print(n)
  fmt.Println()
  return c
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


func walkDocs(prog *loader.Program, t types.Type) {

  switch t := t.(type) {
  case *types.Named:
    u := t.Underlying()
    walkDocs(prog, u)

  case *types.Struct:
    for i := 0; i < t.NumFields(); i++ {
      f := t.Field(i)
      fmt.Println("FIELD", f.Name(), f.Type())

      _, path, _ := prog.PathEnclosingInterval(f.Pos(), f.Pos())
      for _, n := range path {
        d := extractFieldDoc(n)
        if d != "" {
          fmt.Println("DOC", f.Name(), f.Id(), d)
        }
      }

      walkDocs(prog, f.Type())
    }
  case *types.Basic:
  default:
    fmt.Println("unknown type", t)
  }
}

func ParseComments() {

  var conf loader.Config
  _, err := conf.FromArgs([]string{"github.com/buchanae/roger/example/server"}, false)
  conf.ParserMode = parser.ParseComments

	if err != nil {
		panic(err)
	}

  prog, err := conf.Load()
	if err != nil {
		panic(err)
	}

  pkg := prog.Package("github.com/buchanae/roger/example/server")
  co := pkg.Pkg.Scope().Lookup("Config")
  walkDocs(prog, co.Type())

  /*
  o, p, q := types.LookupFieldOrMethod(st, true, pkg.Pkg, "HostName")
  fmt.Println(o, p, q)
  */

  /*
  for e, tv := range pkg.Types {
    fmt.Println(e, tv)
  }
  ?

  /*
  for _, f := range pkg.Files {
    pretty.Print(f)
  }

  //c := collector{}
  files := map[string]*ast.File{}

  for _, f := range pkg.Files {
    tokfile := prog.Fset.File(f.Pos())
    name := tokfile.Name()
    files[name] = f
  }

  astpkg := ast.Package{
    Name: "server",
    Files: files,
  }

  d := doc.New(&astpkg, "github.com/buchanae/roger/example/server", doc.AllDecls)
  for _, t := range d.Types {
    if t.Name == "Config" {
      //ast.Walk(c, t.Decl)
      for _, s := range t.Decl.Specs {
        if ts, ok := s.(*ast.TypeSpec); ok {
          if st, ok := ts.Type.(*ast.StructType); ok {
            for _, f := range st.Fields.List {
              name := f.Names[0].Name
              var text []string
              if f.Doc != nil {
                for _, d := range f.Doc.List {
                  text = append(text, d.Text)
                }
              }
              pretty.Print(f)
              fmt.Println()
              fmt.Println(name, text)
            }
          }
        }
      }
    }
  }
  */

  /*
  for _, pkginfo := range prog.Imported {
    for _, f := range pkginfo.Files {

      ast.Inspect(f, func(n ast.Node) bool {
        if t, ok := n.(*ast.TypeSpec); ok {
          if t.Name.Name == "Config" {
            if s, ok := t.Type.(*ast.StructType); ok {
              for _, f := range s.Fields.List {
                fmt.Println(f)
                fmt.Println(f.Doc.Text())
                fmt.Println(f.Names)
              }
              c := collector{}
              ast.Walk(c, s)
            }
          }
        }
        return true
      })
    }
  }
  */
}

func TestRoger(t *testing.T) {
  ParseComments()
}

func DontTestRoger(t *testing.T) {

  c := DefaultConfig()
  tr := Inspect(&c, []string{
    "scheduler.worker",
  })
  fs := flag.NewFlagSet("roger", flag.ExitOnError)

  alias := map[string]string{
    "server.host_name": "host",
    "worker.work_dir": "w",
  }

  for name, l := range tr.leaves {
    fmt.Printf("%-60s %s\n", name, l.Type)
    fs.Var(l, name, "usage")

    if a, ok := alias[name]; ok {
      fs.Var(l, a, "usage")
    }
  }

  yamlconf, err := loadYAML("default-config.yaml")
  if err != nil {
    fmt.Println(err)
  }

  yamlflat := map[string]interface{}{}
  flatten(yamlconf, "", yamlflat)
  setValues(tr.leaves, yamlflat)

  tr.LoadEnv(envname)

  args := []string{
    "-worker.active_event_writers", "baz",
    "-worker.active_event_writers", "bat lak",
    "-w", "flagsetworkdir",
    //"-worker.task_reader", "foo",
    //"-worker.update_rate", "20s",

    // invalid
    //"-scheduler.schedule_chunk", "z",
  }

  err = fs.Parse(args)
  if err != nil {
    fmt.Println(err)
  }
  //fs.PrintDefaults()

  //tr.ignoreEmpty = true

  c.Scheduler.Worker = c.Worker
  c.Log.Level = "blah"

  tr.dump(nil)
  fmt.Println()
  fmt.Println(c.Worker.ActiveEventWriters)
  fmt.Println(c.Worker.WorkDir)
  fmt.Println(c.Worker.Storage.Local.AllowedDirs)
  fmt.Println(c.Worker.UpdateRate)

  for _, err := range tr.validate(nil) {
    fmt.Println(err)
  }

}

func setValues(dest map[string]*leaf, src map[string]interface{}) {
  for name, val := range src {

    // TODO
    // If there's a block defined but all its values are commented out,
    // this will show up as unknown. Debatable what should be done in that case.
    // It isn't technically unknown, but it's not very clean either.
    if val == nil {
      continue
    }

    l, ok := dest[name]
    if !ok {
      fmt.Println("unknown", name)
      continue
    }

    if err := l.Coerce(val); err != nil {
      fmt.Println(err)
    }
  }
}


func join(path []string, delim string, prefix string, transform func(string) string) string {
  var p []string
  if prefix != "" {
    p = append(p, prefix)
  }
  for _, i := range path {
    p = append(p, transform(i))
  }
  return strings.Join(p, delim)
}

func flagname(path []string) string {
  return join(path, ".", "", underscore)
}
func envname(path []string) string {
  return join(path, "_", "funnel", underscore)
}

func flatten(in map[string]interface{}, prefix string, out map[string]interface{}) {
  for k, v := range in {
    path := k
    if prefix != "" {
      path = prefix + "." + k
    }
    path = flagname(strings.Split(path, "."))

    switch x := v.(type) {
    case map[string]interface{}:
      flatten(x, path, out)
    default:
      out[path] = v
    }
  }
}



func loadJSON(path string) (map[string]interface{}, error) {
  jsonconf := map[string]interface{}{}
  jsonb, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }

  err = json.Unmarshal(jsonb, &jsonconf)
  if err != nil {
    return nil, err
  }
  return jsonconf, nil
}

func loadYAML(path string) (map[string]interface{}, error) {
  yamlconf := map[string]interface{}{}
  yamlb, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }

  err = yaml.Unmarshal(yamlb, &yamlconf)
  if err != nil {
    return nil, err
  }
  return yamlconf, nil
}


func newpathI(base []int, add ...int) []int {
  path := append([]int{}, base...)
  return append(path, add...)
}
func newpathS(base []string, add ...string) []string {
  path := append([]string{}, base...)
  return append(path, add...)
}
