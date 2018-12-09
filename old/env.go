package roger

import (
	"os"
)

// EnvProvider provides access to config values from environment variables.
type EnvProvider struct {
  Prefix string
}

// NewEnvProvider returns a new EnvProvider for accessing environment variables
// under the given "prefix", e.g. "prefix_root_sub_sub_one".
func Env(prefix string) *EnvProvider {
	return &EnvProvider{Prefix: prefix}
}

// Init initialized the provider.
func (e *EnvProvider) Init() error {
	return nil
}

// Lookup returns the value of the given "key", where "key" looks like "Root.Sub.SubOne".
// If the key isn't found in the environment, nil is returned.
func (e *EnvProvider) Lookup(key string) (interface{}, error) {
  var envKeyFunc = PrefixKey(prefix, UnderscoreKey)
	key = tryKeyfunc(key, envKeyFunc, UnderscoreKey)
	v, ok := os.LookupEnv(key)
	if !ok {
		return nil, nil
	}
	return v, nil
}
