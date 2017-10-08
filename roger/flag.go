package roger

import (
  "flag"
  "strings"
  "os"
  "github.com/alecthomas/units"
  "time"
)

type FlagProvider struct {
  Keyfunc
  *flag.FlagSet
}

func NewFlagProvider(fs *flag.FlagSet) *FlagProvider {
  return &FlagProvider{FlagSet: fs}
}

func (f *FlagProvider) Init() error {
  if !f.FlagSet.Parsed() {
    return f.FlagSet.Parse(os.Args[1:])
  }
  return nil
}

func (f *FlagProvider) Lookup(key string) (interface{}, error) {
  key = tryKeyfunc(key, f.Keyfunc, FlagKey)
  fl := f.FlagSet.Lookup(key)
  if fl == nil {
    return nil, nil
  }
  return fl.Value.(flag.Getter).Get(), nil
}

func FlagKey(k string) string {
  return join(strings.Split(k, "."), ".", "", underscore)
}

func AddFlags(rv RogerVals, fs *flag.FlagSet, kf Keyfunc) {
  for k, v := range rv.RogerVals() {
    k = tryKeyfunc(k, kf, FlagKey)
    switch x := v.val.(type) {
    case *uint:
      fs.Uint(k, *x, v.Doc)
    case *uint64:
      fs.Uint64(k, *x, v.Doc)
    case *int:
      fs.Int(k, *x, v.Doc)
    case *int64:
      fs.Int64(k, *x, v.Doc)
    case *float64:
      fs.Float64(k, *x, v.Doc)
    case *bool:
      fs.Bool(k, *x, v.Doc)
    case *string:
      fs.String(k, *x, v.Doc)
    case *time.Duration:
      fs.Duration(k, *x, v.Doc)
    case *[]string, *units.MetricBytes:
      fs.Var(&flagVal{}, k, v.Doc)
    }
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
