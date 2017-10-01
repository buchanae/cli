package roger

import (
  "fmt"
  "flag"
  "testing"
  "os"
  "io/ioutil"
  "encoding/json"
	"github.com/ghodss/yaml"
  "github.com/buchanae/roger/example/server"
  "github.com/buchanae/roger/example/worker"
  "github.com/buchanae/roger/example/scheduler"
  "github.com/buchanae/roger/example/logger"
  "github.com/buchanae/roger/example/dynamo"
  "github.com/spf13/cast"
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

TODO:
- dump flag, env, yaml, json
- alias/link/source field value from another field
- ignore/hide fields
- define short fields
- validation interface
- support byte size
- support SI prefix (K, G, M, etc)
- choices + validation
- printing config, but only non-defaults
- help/docs from comments
- manage editing config file
- pluggable sources
- sets of default configurations
- slice of choices
- improve stringSlice.String() format
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

func Inspect(i interface{}) *tree {
  t := reflect.TypeOf(i)
  v := reflect.ValueOf(i)

  if v.Kind() != reflect.Ptr || v.IsNil() {
    panic("must be non-nil pointer type")
  }

  // TODO check that it's a struct type

  tr := tree{
    st: t.Elem(),
    sv: v.Elem(),
  }
  tr.inspect(nil)
  return &tr
}

type tree struct {
  leaves []*leaf
  st reflect.Type
  sv reflect.Value
  ignoreEmpty bool
}

type leaf struct {
  Path []string
  Type reflect.Type
  Value reflect.Value
  Addr interface{}
}

func (tr *tree) pathname(path []int) []string {
  var name []string
  for i := 0; i < len(path); i++ {
    name = append(name, tr.st.FieldByIndex(path[:i+1]).Name)
  }
  return name
}

func (tr *tree) dump(base []int) {
  t := tr.sv.FieldByIndex(base)

  for j := 0; j < t.NumField(); j++ {
    indent := strings.Repeat("  ", len(base))
    path := newpathI(base, j)
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

func (tr *tree) inspect(base []int) {
  t := tr.sv.FieldByIndex(base)

  for j := 0; j < t.NumField(); j++ {
    path := newpathI(base, j)

    fv := tr.sv.FieldByIndex(path)

    //indent := strings.Repeat("  ", len(path))

    switch fv.Kind() {
    case reflect.Struct:
      tr.inspect(path)

    default:
      //fmt.Println(indent, ft.Name, ":", fv)
      tr.leaves = append(tr.leaves, &leaf{
        Path: tr.pathname(path),
        Type: fv.Type(),
        Value: fv,
        Addr: fv.Addr().Interface(),
      })
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

func TestRoger(t *testing.T) {

  c := DefaultConfig()
  tr := Inspect(&c)
  res := tr.leaves
  fs := flag.NewFlagSet("roger", flag.ExitOnError)
  byname := map[string]*leaf{}

  for _, k := range res {
    name := flagname(k.Path)
    byname[name] = k

    fmt.Printf("%-60s %s\n", name, k.Type)

    switch x := k.Value.Interface().(type) {
    case string:
      fs.StringVar(k.Addr.(*string), name, x, "usage")

    case int:
      fs.IntVar(k.Addr.(*int), name, x, "usage")

    case int64:
      fs.Int64Var(k.Addr.(*int64), name, x, "usage")

    case bool:
      fs.BoolVar(k.Addr.(*bool), name, x, "usage")

    case float64:
      fs.Float64Var(k.Addr.(*float64), name, x, "usage")

    case uint:
      fs.UintVar(k.Addr.(*uint), name, x, "usage")

    case uint64:
      fs.Uint64Var(k.Addr.(*uint64), name, x, "usage")

    case []string:
      fs.Var(sliceVar{dest: k.Addr.(*[]string)}, name, "usage")

    case time.Duration:
      fs.DurationVar(k.Addr.(*time.Duration), name, x, "usage")
    }
  }

  yamlconf, err := loadYAML("default-config.yaml")
  if err != nil {
    fmt.Println(err)
  }

  jsonconf, err := loadJSON("default-config.json")
  if err != nil {
    fmt.Println(err)
  }

  yamlflat := map[string]interface{}{}
  flatten(yamlconf, "", yamlflat)

  jsonflat := map[string]interface{}{}
  flatten(jsonconf, "", jsonflat)

  for k, v := range yamlflat {
    fmt.Println(flagname(strings.Split(k, ".")), v)
  }

  fmt.Println()

  for k, v := range jsonflat {
    fmt.Println(flagname(strings.Split(k, ".")), v)
  }

  setValues(byname, yamlflat)
  setValues(byname, jsonflat)

  var envargs []string
  for _, k := range res {
    v := os.Getenv(envname(k.Path))
    if v != "" {
      envargs = append(envargs, "-" + flagname(k.Path), v)
    }
  }

  err = fs.Parse(envargs)
  if err != nil {
    fmt.Println(err)
  }

  args := []string{
    "-worker.active_event_writers", "baz",
    "-worker.work_dir", "flagsetworkdir",
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

  fmt.Println(c.Worker.ActiveEventWriters)
  fmt.Println(c.Worker.WorkDir)
  fmt.Println(c.Worker.Storage.Local.AllowedDirs)
  fmt.Println(c.Worker.UpdateRate)

  tr.ignoreEmpty = true
  tr.dump(nil)
}

func setValues(dest map[string]*leaf, src map[string]interface{}) {
  for name, val := range src {
    fmt.Println("SETTING", name, val)

    // TODO
    // If there's a block defined but all its values are commented out,
    // this will show up as unknown. Debatable what should be done in that case.
    // It isn't technically unknown, but it's not very clean either.
    if val == nil {
      continue
    }

    l, ok := dest[name]
    if !ok {
      //fmt.Println("unknown", name)
      continue
    }


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
    case time.Duration:
      casted, err = cast.ToDurationE(val)
    default:
      fmt.Println("unknown source value", name, val)
      continue
    }

    if err != nil {
      fmt.Println("error casting", name, val, err)
      continue
    }

    l.Value.Set(reflect.ValueOf(casted))
  }
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



type sliceVar struct {
  dest *[]string
  cleared bool
}
func (sv sliceVar) String() string {
  if sv.dest == nil {
    sv.dest = &[]string{}
  }
  return strings.Join(*sv.dest, " ")
}
func (sv sliceVar) Set(s string) error {
  if sv.dest == nil {
    sv.dest = &[]string{}
  }
  if !sv.cleared {
    sv.cleared = true
    *sv.dest = []string{}
  }
  *sv.dest = append(*sv.dest, s)
  return nil
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
