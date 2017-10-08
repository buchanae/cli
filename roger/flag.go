package roger

import (
	"flag"
	"github.com/alecthomas/units"
	"os"
	"strings"
	"time"
)

// FlagProvider provides access to values in a flag.FlagSet.
//
// By default, flag keys look like "root.sub.sub_one"
type FlagProvider struct {
	Keyfunc
	Flags *flag.FlagSet
}

// NewFlagProvider returns a new FlagProvider for the given RogerVals.
func NewFlagProvider(rv RogerVals) *FlagProvider {
	fp := &FlagProvider{Flags: &flag.FlagSet{}}
	fp.AddFlags(rv)
	return fp
}

// Init will call FlagSet.Parse() if it has not been called yet.
func (f *FlagProvider) Init() error {
	if !f.Flags.Parsed() {
		return f.Flags.Parse(os.Args[1:])
	}
	return nil
}

// Lookup returns the flag value of the given key, if it was set,
// where "key" looks like "Root.Sub.SubOne".
func (f *FlagProvider) Lookup(key string) (interface{}, error) {
	key = tryKeyfunc(key, f.Keyfunc, FlagKey)
	fl := f.Flags.Lookup(key)
	if fl == nil {
		return nil, nil
	}
	return fl.Value.(flag.Getter).Get(), nil
}

// FlagKey is the default Keyfunc for flag values.
// "Root.Sub.SubOne" is transformed to "root.sub.sub_one"
func FlagKey(k string) string {
	return join(strings.Split(k, "."), ".", underscore)
}

// AddFlags adds flags for all the configurable keys in "rv".
func (f *FlagProvider) AddFlags(rv RogerVals) {
	fs := f.Flags
	for k, v := range rv.RogerVals() {
		k = tryKeyfunc(k, f.Keyfunc, FlagKey)
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
