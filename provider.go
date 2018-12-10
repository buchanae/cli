package cli

type Provider interface {
	Init() error
	Lookup(key []string) (interface{}, bool)
}

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

func LoadOpts(p Provider, opts []OptSpec) error {
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
