package cli

// Provider is implemented by types which provide config values,
// such as from environment variables, config files, CLI flags, etc.
type Provider interface {
	Init() error
  // Lookup looks up a value for the given key. 
  // "exists" is true if the value was set.
	Lookup(key []string) (value interface{}, exists bool)
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

// LoadOpts loads and sets option values from the given Provider.
func LoadOpts(opts []OptSpec, p Provider) error {
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
