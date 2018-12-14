package cli

import (
	"github.com/spf13/pflag"
)

// PFlags loads option values from a pflag.FlagSet.
type pflags struct {
	KeyFunc
	*pflag.FlagSet
	flags []*pflagValue
}

func PFlags(fs *pflag.FlagSet, opts []*Opt, kf KeyFunc) Provider {
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

func (f *pflags) Provide(l *Loader) error {
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
	p.val = v
	p.set = true
	return nil
}

func (p *pflagValue) String() string {
	return p.opt.DefaultString
}

func (p *pflagValue) Type() string {
	// TODO shows io.Writer, which doesn't make much sense for CLI.
	return p.opt.Type
}
