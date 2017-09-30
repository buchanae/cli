package roger

import (
  "fmt"
  "flag"
  "testing"
  "os"
  "github.com/buchanae/roger/example/server"
  "github.com/buchanae/roger/example/worker"
  "github.com/buchanae/roger/example/scheduler"
  "github.com/buchanae/roger/example/logger"
  "github.com/buchanae/roger/example/dynamo"
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
- read from flag, env, yaml, json
- dump flag, env, yaml, json
- manage editing config file
- alias/link/source field value from another field
- ignore/hide fields
- define short fields
- report unknown fields
- help/docs from comments
- pluggable sources
- support time.Duration in yaml, json, env
- case sensitivity
- choices + validation
- printing config, but only non-defaults
- sets of default configurations
- validation interface
- support byte size
- support SI prefix (K, G, M, etc)
- slice of choices
- improve stringSlice.String() format
- handle map[string]string via "key=value" flag value
- explore "storage.local.allowed_dirs.append"

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
  Value reflect.Value
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
        Value: fv,
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

func TestRoger(t *testing.T) {

  c := DefaultConfig()
  res := Inspect(&c)

  flagname := func(path []string) string {
    return join(path, ".", "", underscore)
  }
  envname := func(path []string) string {
    return join(path, "_", "funnel", underscore)
  }

  for _, k := range res {
    fmt.Printf("%-60s %s\n",
      flagname(k.Path),
      reflect.TypeOf(k.Value.Interface()),
    )
  }

  fs := flag.NewFlagSet("roger", flag.ExitOnError)
  for _, k := range res {
    name := flagname(k.Path)
    p := k.Value.Addr().Interface()

    switch x := k.Value.Interface().(type) {

    case string:
      fs.StringVar(p.(*string), name, x, "usage")

    case int:
      fs.IntVar(p.(*int), name, x, "usage")

    case int64:
      fs.Int64Var(p.(*int64), name, x, "usage")

    case bool:
      fs.BoolVar(p.(*bool), name, x, "usage")

    case float64:
      fs.Float64Var(p.(*float64), name, x, "usage")

    case uint:
      fs.UintVar(p.(*uint), name, x, "usage")

    case uint64:
      fs.Uint64Var(p.(*uint64), name, x, "usage")

    case []string:
      fs.Var(sliceVar{dest: p.(*[]string)}, name, "usage")

    case time.Duration:
      fs.DurationVar(p.(*time.Duration), name, x, "usage")
    }
  }

  args := []string{}
  //args := []string{"-worker.task_reader", "foo"}
  args = []string{"-worker.active_event_writers", "baz", "-worker.work_dir", "flagset"}

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

  fmt.Println(c.Worker.ActiveEventWriters)
  fmt.Println(c.Worker.WorkDir)
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
