package roger

import (
  "fmt"
  "os"
  "flag"
  "strings"
  "path/filepath"
  "reflect"
)

type Vals map[string]Val

func (v Vals) DeletePrefix(prefix ...string) {
  for k, _ := range v {
    for _, p := range prefix {
      if strings.HasPrefix(k, p) {
        delete(v, k)
      }
    }
  }
}

func (v Vals) Alias(a map[string]string) {
  for to, from := range a {
    if _, ok := v[from]; ok {
      v[to] = v[from]
    }
  }
}

type Validator interface {
  Validate() []error
}

type Keyfunc func(string) string

type Val struct {
  Doc string
  val interface{}
}

func NewVal(doc string, v interface{}) Val {
  return Val{doc, v}
}

func FromEnv(v Val, key string) error {
  env, ok := os.LookupEnv(key)
  if ok {
    return CoerceSet(v, env)
  }
  return nil
}

func FromAllEnvPrefix(vals Vals, prefix string) {
  FromAllEnvFunc(vals, PrefixEnvKey(prefix))
}

func FromAllEnv(vals Vals) {
  for k, v := range vals {
    FromEnv(v, EnvKey(k))
  }
}

func FromAllEnvFunc(vals Vals, kf Keyfunc) {
  for k, v := range vals {
    FromEnv(v, kf(k))
  }
}

func AddFlags(vals Vals, fs *flag.FlagSet) {
  for k, v := range vals {
    fv := &FlagVal{val: v}
    fs.Var(fv, FlagKey(k), v.Doc)
  }
}

func AddFlagsFunc(vals Vals, fs *flag.FlagSet, kf Keyfunc) {
  for k, v := range vals {
    fs.Var(&FlagVal{val: v}, kf(k), v.Doc)
  }
}


type UnknownField string
func (u UnknownField) Error() string {
  return fmt.Sprintf("unknown field: %s", string(u))
}

func IsUnknownField(err error) (string, bool) {
  if f, ok := err.(UnknownField); ok {
    return string(f), true
  }
  return "", false
}

func FromMap(vals Vals, m map[string]interface{}) []error {
  var errs []error
  f := map[string]interface{}{}
  flatten(m, "", f)

  for fk, fv := range f {
    // TODO
    // If there's a block defined but all its values are commented out,
    // this will show up as unknown. Debatable what should be done in that case.
    // It isn't technically unknown, but it's not very clean either.
    if fv == nil {
      continue
    }

    rv, ok := vals[fk]
    if !ok {
      errs = append(errs, UnknownField(fk))
      continue
    }

    if err := CoerceSet(rv, fv); err != nil {
      fmt.Println(err)
    }
  }
  return errs
}

func flatten(in map[string]interface{}, prefix string, out map[string]interface{}) {
  for k, v := range in {
    path := k
    if prefix != "" {
      path = prefix + "." + k
    }

    switch x := v.(type) {
    case map[string]interface{}:
      flatten(x, path, out)
    default:
      out[path] = v
    }
  }
}

func FromFile(vals Vals, path string) []error {
  ext := filepath.Ext(path)
  switch ext {
  case ".yaml", ".yml":
    return FromYAMLFile(vals, path)
  default:
    return []error{fmt.Errorf("unknown file type: %s", ext)}
  }
  return nil
}

func FromYAMLFile(vals Vals, path string) []error {
  y, err := LoadYAML(path)
  if err != nil {
    return []error{err}
  }
  return FromMap(vals, y)
}

type FlagVal struct {
  val Val
}

func (f *FlagVal) String() string {
  // TODO val is a pointer, but maybe want to deref?
  return "TODO val"
}

func (f *FlagVal) Set(s string) error {
  return CoerceSet(f.val, s)
}

func (f *FlagVal) Get() interface{} {
  // TODO returning a pointer, but maybe want to deref?
  return f.val.val
}




func FlagKey(k string) string {
  return join(strings.Split(k, "."), ".", "", underscore)
}

func EnvKey(k string) string {
  return join(strings.Split(k, "."), "_", "", underscore)
}

func PrefixEnvKey(prefix string) Keyfunc {
  return func(k string) string {
    return EnvKey(prefix + "." + k)
  }
}

func Dump(vals Vals) {
  for k, j := range vals {
    fmt.Println(j.val, FlagKey(k))
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

func Validate(i interface{}, ignore []string) []error {
  v := validator{
    rootType: reflect.TypeOf(i).Elem(),
    rootVal: reflect.ValueOf(i).Elem(),
    ignore: map[string]struct{}{},
  }
  for _, i := range ignore {
    v.ignore[i] = struct{}{}
  }
  v.validate(nil)
  return v.errors
}

type validator struct {
  rootType reflect.Type
  rootVal reflect.Value
  ignore map[string]struct{}
  errors []error
}

func (v *validator) validate(base []int) {
  t := v.rootVal.FieldByIndex(base)

  for j := 0; j < t.NumField(); j++ {
    path := newpathI(base, j)
    name := pathname(v.rootType, path)

    if _, ok := v.ignore[name]; ok {
      continue
    }

    ft := v.rootType.FieldByIndex(path)
    fv := v.rootVal.FieldByIndex(path)

    // Ignore unexported fields.
    if ft.PkgPath != "" {
      continue
    }

    if x, ok := fv.Interface().(Validator); ok {
      for _, err := range x.Validate() {
        v.errors = append(v.errors, fmt.Errorf("%s: %s", name, err))
      }
    }

    if fv.Kind() == reflect.Struct {
      v.validate(path)
    }
  }
}

func pathname(root reflect.Type, path []int) string {
  var name []string
  for i := 0; i < len(path); i++ {
    name = append(name, root.FieldByIndex(path[:i+1]).Name)
  }
  return strings.Join(name, ".")
}

func newpathI(base []int, add ...int) []int {
  path := append([]int{}, base...)
  return append(path, add...)
}
