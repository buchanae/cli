package roger

import (
  "fmt"
  "os"
  "flag"
  "strings"
)

type Vals map[string]Val

func (v Vals) DeletePrefix(prefix string) {
  for k, _ := range v {
    if strings.HasPrefix(k, prefix) {
      delete(v, k)
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

func SetFromEnv(v Val, key string) error {
  env, ok := os.LookupEnv(key)
  if ok {
    return CoerceSet(v, env)
  }
  return nil
}

func SetAllFromEnvPrefix(vals Vals, prefix string) {
  SetAllFromEnvFunc(vals, PrefixEnvKey(prefix))
}

func SetAllFromEnv(vals Vals) {
  for k, v := range vals {
    SetFromEnv(v, EnvKey(k))
  }
}

func SetAllFromEnvFunc(vals Vals, kf Keyfunc) {
  for k, v := range vals {
    SetFromEnv(v, kf(k))
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

func SetFromMap(vals Vals, m map[string]interface{}) {
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
      fmt.Println("unknown", fk)
      continue
    }

    if err := CoerceSet(rv, fv); err != nil {
      fmt.Println(err)
    }
  }
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
