package roger

import (
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

type Keyfunc func([]string) string

type Val struct {
  Key []string
  Doc string
  val interface{}
}

func NewVal(k []string, doc string, v interface{}) Val {
  return Val{k, doc, v}
}

func SetFromEnv(v Val, key string) error {
  // TODO allow empty value to unset?
  env, _ := os.LookupEnv(key)
  return CoerceSet(v.val, env)
}

func SetAllFromEnv(vals Vals) {
  for _, v := range vals.RogerVals() {
    SetFromEnv(v, EnvKey(v.Key))
  }
}

func SetAllFromEnvFunc(vals Vals, kf Keyfunc) {
  for _, v := range vals.RogerVals() {
    SetFromEnv(v, kf(v.Key))
  }
}

func AddFlags(fs *flag.FlagSet, vals Vals) {
  for k, v := range vals.RogerVals() {
    fs.Var(FlagVal{val: v, k: k}, FlagKey(v.Key), v.Doc)
  }
}

func AddFlagsFunc(fs *flag.FlagSet, vals Vals, kf Keyfunc) {
  for k, v := range vals.RogerVals() {
    fs.Var(FlagVal{val: v, k: k}, kf(v.Key), v.Doc)
  }
}

type FlagVal struct {
  val Val
  k string
}

func (f FlagVal) String() string {
  // TODO val is a pointer, but maybe want to deref?
  return "TODO val"
}

func (f FlagVal) Set(s string) error {
  return CoerceSet(f.val.val, s)
}

func (f FlagVal) Get() interface{} {
  // TODO returning a pointer, but maybe want to deref?
  return f.val.val
}




func FlagKey(k []string) string {
  return join(k, ".", "", underscore)
}

func EnvKey(k []string) string {
  return join(k, "_", "", underscore)
}

func PrefixEnvKey(prefix string) Keyfunc {
  return func(k []string) string {
    i := append([]string{}, prefix)
    i = append(i, k...)
    return EnvKey(i)
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
