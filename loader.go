package cli

import (
	"fmt"
	"strings"
)

// Provider is implemented by types which provide option values,
// such as from environment variables, config files, CLI flags, etc.
type Provider interface {
	Provide(*Loader) error
}

// NewLoader returns a Loader instance which is configured
// to load option values from the given providers.
func NewLoader(opts []*Opt, providers ...Provider) *Loader {
	var keys [][]string
	for _, opt := range opts {
		keys = append(keys, opt.Key)
	}
	return &Loader{
		keys:      keys,
		opts:      opts,
		providers: providers,
		Coerce:    Coerce,
	}
}

// Loader is used to load, coerce, and set option values
// at command run time. Loader.Load runs the providers
// in order. Once an option has been set, it is not overridden
// by later values.
type Loader struct {
	opts      []*Opt
	keys      [][]string
	providers []Provider
	errors    []error
	// Coerce can be used to override the type coercion
	// needed when setting an option value. A coerce function
	// must set the value. "dst" is always a pointer to the
	// value which needs to be set, e.g. *int for an option
	// value of type int. See coerce.go for an example.
	Coerce func(dst, src interface{}) error
}

// Load runs the providers, loading and setting option values.
// Load is meant to be called only once.
// TODO this is awkward. The code should make it clear that
//      it only gets called once, and/or make it safe to load
//      multiple times.
func (l *Loader) Load() {
	for _, src := range l.providers {
		err := src.Provide(l)
		if err != nil {
			l.errors = append(l.errors, err)
		}
	}
}

// Errors returns a list of errors encountered during loading.
func (l *Loader) Errors() []error {
	return l.errors
}

// Keys returns a list of keys for all options.
func (l *Loader) Keys() [][]string {
	return l.keys
}

// Get gets the current option value for the given key.
func (l *Loader) Get(key []string) interface{} {
	for _, opt := range l.opts {
		if l.eq(key, opt.Key) {
			return opt.Value
		}
	}
	return nil
}

// GetString returns the option value as a string,
// or else an empty string.
func (l *Loader) GetString(key []string) string {
	val := l.Get(key)
	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}

// Set sets an option value for the option at the given key.
// Once an option is set, it will not be overridden by future
// calls to Set. Set uses Loader.Coerce to set the value.
//
// TODO this is ugly. This checks IsSet, but doesn't set it to true,
//      that happens somewhere else. The whole set/coerce interface
//      needs cleanup
func (l *Loader) Set(key []string, val interface{}) {
	for _, opt := range l.opts {
		if !l.eq(key, opt.Key) {
			continue
		}
		if opt.IsSet {
			continue
		}
		err := l.Coerce(opt.Value, val)
		if err != nil {
			l.errors = append(l.errors, err)
		}
		return
	}
	// TODO these errors are missing context, e.g. "in file config.yaml"
	l.errors = append(l.errors, fmt.Errorf("unknown opt key %v", key))
}

// eq returns true if two option keys are equal.
// eq is case insensitive.
func (l *Loader) eq(key1, key2 []string) bool {
	if len(key1) != len(key2) {
		return false
	}
	for i := 0; i < len(key1); i++ {
		if strings.ToLower(key1[i]) != strings.ToLower(key2[i]) {
			return false
		}
	}
	return true
}
