package roger

import (
  "flag"
  "strings"
)

type FlagProvider struct {
  Keyfunc
  *flag.FlagSet
}

func NewFlagProvider(fs *flag.FlagSet) *FlagProvider {
  return &FlagProvider{FlagSet: fs}
}

func (f *FlagProvider) Lookup(key string) (interface{}, bool) {
  key = tryKeyfunc(key, f.Keyfunc, FlagKey)
  fl := f.FlagSet.Lookup(key)
  if fl == nil {
    return nil, false
  }
  val := fl.Value.(flag.Getter).Get()
  if val == nil {
    return nil, false
  }
  return val, true
}

func FlagKey(k string) string {
  return join(strings.Split(k, "."), ".", "", underscore)
}

func AddFlags(rv RogerVals, fs *flag.FlagSet, kf Keyfunc) {
  for k, v := range rv.RogerVals() {
    k = tryKeyfunc(k, kf, FlagKey)
    fs.Var(&flagVal{}, k, v.Doc)
  }
}

type flagVal struct {
  val string
  set bool
}

func (f *flagVal) String() string {
  return f.val
}

func (f *flagVal) Set(s string) error {
  f.val = s
  f.set = true
  return nil
}

func (f *flagVal) Get() interface{} {
  if !f.set {
    return nil
  }
  return f.val
}
