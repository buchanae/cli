package roger

import (
	"os"
)

// EnvProvider provides access to config values from environment variables.
// By default, value keys look like "prefix_root_sub_sub_one".
type EnvProvider struct {
	Keyfunc
}

// NewEnvProvider returns a new EnvProvider for accessing environment variables
// under the given "prefix", e.g. "prefix_root_sub_sub_one".
func NewEnvProvider(prefix string) *EnvProvider {
	return &EnvProvider{
    Keyfunc: PrefixKey(prefix, UnderscoreKey),
  }
}

// Init initialized the provider.
func (e *EnvProvider) Init() error {
	return nil
}

// Lookup returns the value of the given "key", where "key" looks like "Root.Sub.SubOne".
// If the key isn't found in the environment, nil is returned.
func (e *EnvProvider) Lookup(key string) (interface{}, error) {
	key = tryKeyfunc(key, e.Keyfunc, UnderscoreKey)
	v, ok := os.LookupEnv(key)
	if !ok {
		return nil, nil
	}
	return v, nil
}
