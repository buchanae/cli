package roger

import (
  "fmt"
  "flag"
  "testing"
  "os"
  "io/ioutil"
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

TODO:
- read from yaml, json
- dump flag, env, yaml, json
- alias/link/source field value from another field
- ignore/hide fields
- define short fields
- validation interface
- support time.Duration in yaml, json, env
- support byte size
- support SI prefix (K, G, M, etc)
- case sensitivity
- choices + validation
- report unknown fields
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

func Inspect(i interface{}) []*leaf {
  t := reflect.TypeOf(i)
  v := reflect.ValueOf(i)

  if v.Kind() != reflect.Ptr || v.IsNil() {
    panic("must be non-nil pointer type")
  }

  // TODO check that it's a struct type

  is := inspectState{
    st: t.Elem(),
    sv: v.Elem(),
  }
  is.inspect(nil)
  return is.res
}

type inspectState struct {
  res []*leaf
  st reflect.Type
  sv reflect.Value
}

type leaf struct {
  Path []string
  Type reflect.Type
  Value reflect.Value
  Addr interface{}
}

func (is *inspectState) pathname(path []int) []string {
  var name []string
  for i := 0; i < len(path); i++ {
    name = append(name, is.st.FieldByIndex(path[:i+1]).Name)
  }
  return name
}

func (is *inspectState) inspect(path []int) {
  t := is.sv.FieldByIndex(path)

  for j := 0; j < t.NumField(); j++ {
    index := append([]int{}, path...)
    index = append(index, j)

    ft := is.st.FieldByIndex(index)
    fv := is.sv.FieldByIndex(index)

    indent := strings.Repeat("  ", len(path))

    switch fv.Kind() {
    case reflect.Struct:
      fmt.Println(indent, ft.Name)
      is.inspect(index)

    default:
      fmt.Println(indent, ft.Name, ":", fv)
      is.res = append(is.res, &leaf{
        Path: is.pathname(index),
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
  res := Inspect(&c)
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

  var args []string
  //args := []string{"-worker.task_reader", "foo"}
  args = []string{"-worker.active_event_writers", "baz", "-worker.work_dir", "flagset"}
  //args = []string{"-scheduler.schedule_chunk", "z"}

  err := fs.Parse(args)
  if err != nil {
    fmt.Println(err)
  }
  fs.PrintDefaults()

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

  yamlconf := map[string]interface{}{}
  yamlb, err := ioutil.ReadFile("default-config.yaml")
  if err != nil {
    fmt.Println(err)
  }
  err = yaml.Unmarshal(yamlb, &yamlconf)
  if err != nil {
    fmt.Println(err)
  }

  visit(yamlconf, nil, func(path []string, val interface{}) {
    // TODO
    // If there's a block defined but all its values are commented out,
    // this will show up as unknown. Debatable what should be done in that case.
    // It isn't technically unknown, but it's not very clean either.
    if val == nil {
      return
    }

    name := flagname(path)
    l, ok := byname[name]
    if !ok {
      fmt.Println("unknown", name)
      return
    }

    fmt.Println("setting", name, val)

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
      return
    }

    if err != nil {
      fmt.Println("error casting", name, val, err)
      return
    }

    l.Value.Set(reflect.ValueOf(casted))
  })

  fmt.Println(c.Worker.ActiveEventWriters)
  fmt.Println(c.Worker.WorkDir)
  fmt.Println(c.Worker.Storage.Local.AllowedDirs)
}

type visitor func(path []string, val interface{})
func visit(m map[string]interface{}, base []string, cb visitor) {
  for k, v := range m {

    path := append([]string{}, base...)
    path = append(path, k)

    switch x := v.(type) {
    case map[string]interface{}:
      visit(x, path, cb)
    default:
      cb(path, x)
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
