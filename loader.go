package roger

import (
	"fmt"
)

// Provider defines the interface of config value providers
// such as EnvProvider, FlagProvier, FileProvider, etc.
type Provider interface {
	// Init allows the provider to initialize, and is called by Load().
	// If Init returns an error, Lookup will not be called.
	Init() error
	// Lookup should return the value for the given key.
	// "key" looks like "Root.Sub.SubOne".
	// If the value is not found, the provider should return nil.
	Lookup(key string) (interface{}, error)
}

// Load loads and sets the RogerVals using the given providers.
//
// Init is called on each provider, in order, before any lookups are done.
// If Init returns an error, Lookup will not be called for that provider,
// but loading will continue.
func Load(rv RogerVals, ps ...Provider) []error {
	var errs []error
	var ok []Provider

	for _, p := range ps {
		// Just in case, avoid nil panic
		if p == nil {
			continue
		}

		if err := p.Init(); err != nil {
			errs = append(errs, err)
			continue
		}
		ok = append(ok, p)
	}

	vals := rv.RogerVals()
	for _, p := range ok {
		for k, v := range vals {
			x, err := p.Lookup(k)
			if x != nil && err == nil {
				err = coerceSet(v.val, x)
			}
			if err != nil {
				errs = append(errs, fmt.Errorf("error loading %s: %s", k, err))
			}
		}
	}

	return errs
}
