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

type Loader struct {
	opts      []*Opt
	keys      [][]string
	providers []Provider
	errors    []error
	Coerce    func(dst, src interface{}) error
}

func (l *Loader) Errors() []error {
	return l.errors
}

func (l *Loader) Load() {
	for _, src := range l.providers {
		err := src.Provide(l)
		if err != nil {
			l.errors = append(l.errors, err)
		}
	}
}

func (l *Loader) Keys() [][]string {
	return l.keys
}

// TODO return (val, ok) pattern?
func (l *Loader) Get(key []string) interface{} {
	for _, opt := range l.opts {
		if l.eq(key, opt.Key) {
			return opt.Value
		}
	}
	return nil
}

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

func (l *Loader) GetString(key []string) string {
	val := l.Get(key)
	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}

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
