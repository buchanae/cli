package roger

import (
  "fmt"
  "os"
  "flag"
  "strings"
)

type Vals interface {
  RogerVals() map[string]Val
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
  // TODO allow empty value to unset?
  env, ok := os.LookupEnv(key)
  if ok {
    return CoerceSet(v.val, env)
  }
  return nil
}

func SetAllFromEnvPrefix(vals Vals, prefix string) {
  SetAllFromEnvFunc(vals, PrefixEnvKey(prefix))
}

func SetAllFromEnv(vals Vals) {
  for k, v := range vals.RogerVals() {
    SetFromEnv(v, EnvKey(k))
  }
}

func SetAllFromEnvFunc(vals Vals, kf Keyfunc) {
  for k, v := range vals.RogerVals() {
    SetFromEnv(v, kf(k))
  }
}

func AddFlags(vals Vals, fs *flag.FlagSet) {
  for k, v := range vals.RogerVals() {
    fv := &FlagVal{val: v}
    fs.Var(fv, FlagKey(k), v.Doc)
  }
}

func AddFlagsFunc(vals Vals, fs *flag.FlagSet, kf Keyfunc) {
  for k, v := range vals.RogerVals() {
    fs.Var(&FlagVal{val: v}, kf(k), v.Doc)
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
  return CoerceSet(f.val.val, s)
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

func Dump(v Vals) {
  for k, j := range v.RogerVals() {
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
