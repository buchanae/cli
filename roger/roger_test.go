package roger

import (
  "fmt"
  "testing"
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
- help/docs from comments
- read from flag, env, yaml, json
- support time.Duration in yaml, json, env
- support byte size
- support SI suffix (K, G, M, etc)
- alias/link/source field value from another field
- ignore/hide fields
- define short fields
- report unknown fields
- dump yaml
  - printing config, but only non-defaults

- validation interface
- choices + validation

TODO:
- manage editing config file
- pluggable sources
- slice of choices
- improve stringSlice.String() format
- dump json, env, flags
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
*/

func TestEnvKey(t *testing.T) {
  m := map[string][]string{
    "one_two_three_four": []string{"One", "Two", "ThreeFour"},
  }
  for expected, in := range m {
    got := EnvKey(in)
    if got != expected {
      t.Errorf("expected %s, got %s", expected, got)
    }
  }
}

func TestPrefixEnvKey(t *testing.T) {
  got := PrefixEnvKey("RogerThat")([]string{"One", "TwoThree"})
  if got != "roger_that_one_two_three" {
    t.Errorf("expected roger_that_one_two, got %s", got)
  }
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
    //namer: flagname,
  }
  tr.inspect(nil)
  return &tr
}


type leaf struct {
  Path []string
  Type reflect.Type
  Value reflect.Value
}

type tree struct {
  leaves map[string]*leaf
  hide []string
  st reflect.Type
  sv reflect.Value
  ignoreEmpty bool
  namer func(path []string) string
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
      }
    }
  }
}


/*
func DontTestRoger(t *testing.T) {

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


*/
