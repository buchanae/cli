package cli

import (
  "fmt"
  "strings"
)

// Provider is implemented by types which provide config values,
// such as from environment variables, config files, CLI flags, etc.
type Provider interface {
	Init() error
  // Lookup looks up a value for the given key. 
  // "exists" is true if the value was set.
	Lookup(key []string) (value interface{}, exists bool)
}

// KeysValidator is implemented by Providers when possible
// and allows cli code to report unknown keys.
type KeysValidator interface {
  // ValidateKeys is given a list of allowed keys,
  // and should return a list of errors describing unrecognized keys.
  ValidateKeys(allowed [][]string) []error
}

// Providers returns a Provider that runs the providers in order
// Init stops on the first error; Lookup returns on the first value
// that was found.
func Providers(ps ...Provider) Provider {
	return multiProvider(ps)
}

type multiProvider []Provider

func (m multiProvider) Init() error {
	for _, p := range m {
		err := p.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m multiProvider) Lookup(key []string) (interface{}, bool) {
	for _, p := range m {
		val, ok := p.Lookup(key)
		if ok {
			return val, true
		}
	}
	return nil, false
}

func (m multiProvider) ValidateKeys(allowed [][]string) []error {
  var errs []error
  for _, p := range m {
    vk, ok := p.(KeysValidator)
    if ok {
      errs = append(errs, vk.ValidateKeys(allowed)...)
    }
  }
  return errs
}

func ValidateKeys(opts []OptSpec, p Provider) []error {
  vk, ok := p.(KeysValidator)
  if !ok {
    return nil
  }
  var allowed [][]string
  for _, opt := range opts {
    allowed = append(allowed, opt.Key)
  }
  return vk.ValidateKeys(allowed)
}

// LoadOpts loads and sets option values from the given Provider.
func LoadOpts(opts []OptSpec, p Provider) error {

  err := p.Init()
  if err != nil {
    return err
  }

  errs := ValidateKeys(opts, p)
  if errs != nil {
    var lines []string
    for _, err := range errs {
      lines = append(lines, err.Error())
    }
    return fmt.Errorf(strings.Join(lines, "\n"))
  }

	for _, opt := range opts {
		val, ok := p.Lookup(opt.Key)
		if !ok {
			continue
		}
		err := coerceSet(opt.Value, val)
		if err != nil {
			return err
		}
	}
	return nil
}
