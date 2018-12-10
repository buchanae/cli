package cli

import (
	"github.com/spf13/pflag"
	"os"
	"time"
)

// PFlags returns a Provider that loads values from a pflag.FlagSet.
func PFlags(opts []OptSpec) *PFlagProvider {
	p := &PFlagProvider{}
	p.AddOpts(opts)
	return p
}

// PFlagProvider provides values from a pflag.FlagSet.
type PFlagProvider struct {
	KeyFunc
	*pflag.FlagSet
}

// Init will call FlagSet.Parse() if it has not been called yet.
func (f *PFlagProvider) Init() error {
	if f.FlagSet == nil {
		return nil
	}

	if !f.FlagSet.Parsed() {
		err := f.FlagSet.Parse(os.Args[1:])
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *PFlagProvider) keyfunc(key []string) string {
	if f.KeyFunc != nil {
		return f.KeyFunc(key)
	} else {
		return DotKey(key)
	}
}

// Lookup looks up a value from the FlagSet.
func (f *PFlagProvider) Lookup(key []string) (interface{}, bool) {
	if f.FlagSet == nil {
		return nil, false
	}

	k := f.keyfunc(key)
	fl := f.FlagSet.Lookup(k)
	if fl == nil {
		return nil, false
	}
	// We only want the flags which have been set by the user.
	if !fl.Changed {
		return nil, false
	}
	return fl.Value.String(), true
}

// AddOpts adds flags to the FlagSet for the given options.
func (f *PFlagProvider) AddOpts(opts []OptSpec) {
	if f.FlagSet == nil {
		f.FlagSet = &pflag.FlagSet{}
	}
	fs := f.FlagSet

	for _, opt := range opts {
		det := ParseOptDetail(opt)
		k := f.keyfunc(opt.Key)

		switch z := opt.Value.(type) {
		case *uint:
			fs.UintVar(z, k, *z, det.Synopsis)
		case *uint8:
			fs.Uint8Var(z, k, *z, det.Synopsis)
		case *uint16:
			fs.Uint16Var(z, k, *z, det.Synopsis)
		case *uint32:
			fs.Uint32Var(z, k, *z, det.Synopsis)
		case *uint64:
			fs.Uint64Var(z, k, *z, det.Synopsis)
		case *int:
			fs.IntVar(z, k, *z, det.Synopsis)
		case *int8:
			fs.Int8Var(z, k, *z, det.Synopsis)
		case *int16:
			fs.Int16Var(z, k, *z, det.Synopsis)
		case *int32:
			fs.Int32Var(z, k, *z, det.Synopsis)
		case *int64:
			fs.Int64Var(z, k, *z, det.Synopsis)
		case *float32:
			fs.Float32Var(z, k, *z, det.Synopsis)
		case *float64:
			fs.Float64Var(z, k, *z, det.Synopsis)
		case *bool:
			fs.BoolVar(z, k, *z, det.Synopsis)
		case *string:
			fs.StringVar(z, k, *z, det.Synopsis)
		case *[]string:
			fs.StringSliceVar(z, k, *z, det.Synopsis)
		case *time.Duration:
			fs.DurationVar(z, k, *z, det.Synopsis)
		default:
			// TODO should probably return error
		}
	}
}
