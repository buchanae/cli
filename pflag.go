package cli

import (
  "log"
  "reflect"
  "fmt"
	"github.com/spf13/pflag"
)

// PFlags loads option values from a pflag.FlagSet.
type pflags struct {
	KeyFunc
  *pflag.FlagSet
  flags []*pflagValue
}

func PFlags(fs *pflag.FlagSet, opts []*Opt, kf KeyFunc) Source {
  pf := &pflags{
    KeyFunc: kf,
    FlagSet: fs,
  }

  for _, opt := range opts {
    flag := &pflagValue{opt: opt}
    k := pf.keyfunc(opt.Key)
    fs.VarP(flag, k, opt.Short, opt.Synopsis)

    if opt.Deprecated != "" {
      fs.MarkDeprecated(k, opt.Deprecated)
    }
    if opt.Hidden {
      fs.MarkHidden(k)
    }
    pf.flags = append(pf.flags, flag)
  }
  return pf
}

func (f *pflags) keyfunc(key []string) string {
  if f.KeyFunc == nil {
    return DotKey(key)
  }
  return f.KeyFunc(key)
}

func (f *pflags) Load(l *Loader) error {
  for _, flag := range f.flags {
    if flag.set {
      l.Set(flag.opt.Key, flag.val)
    }
  }
  return nil
}

type pflagValue struct {
  opt *Opt
  val interface{}
  set bool
}

func (p *pflagValue) Set(v string) error {
  log.Println("SET", p.opt.Key, v)
  p.val = v
  p.set = true
  return nil
}

func (p *pflagValue) String() string {
  // TODO not super happy about including reflect just for this.
  //      also not happy with the default doc, want "os.Stdout" for io.Writer
  //      in some cases.
  t := reflect.TypeOf(p.opt.DefaultValue)
  if t.Kind() == reflect.Ptr {
    return ""
  }
  if p.opt.DefaultValue == nil {
    return ""
  }
  return fmt.Sprintf("%v", p.opt.DefaultValue)
}

func (p *pflagValue) Type() string {
  // TODO shows io.Writer, which doesn't make much sense for CLI.
  return p.opt.Type
}
