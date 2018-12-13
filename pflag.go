package cli

import (
	"github.com/spf13/pflag"
	"os"
)

// PFlags loads option values from a pflag.FlagSet.
type PFlags struct {
	KeyFunc
}

func (f *PFlags) Load(l *Loader) error {
  fs := pflag.FlagSet{}

  for _, key := range l.Keys() {
    val := &pflagValue{key, l}
    k := f.KeyFunc(key)
    fs.Var(val, k, "")
  }

  fs.Parse(os.Args[1:])
  return nil
}

type pflagValue struct {
  key []string
  l *Loader
}

func (p *pflagValue) Set(v string) error {
  p.l.Set(p.key, v)
  return nil
}

func (p *pflagValue) String() string {
  return ""
}

func (p *pflagValue) Type() string {
  return ""
}
