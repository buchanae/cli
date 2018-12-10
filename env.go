package cli

import (
	"os"
  "strings"
)

// EnvProvider provides access to config values from environment variables.
type EnvProvider struct {
  KeyFunc
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

func (e *EnvProvider) keyfunc(key []string) string {
  if e.KeyFunc != nil {
    return e.KeyFunc(key)
  } else {
    return strings.ToUpper(UnderscoreKey(key))
  }
}

// Lookup returns the value of the given "key", where "key" looks like "Root.Sub.SubOne".
// If the key isn't found in the environment, nil is returned.
func (e *EnvProvider) Lookup(key []string) (interface{}, bool) {
  key = append([]string{e.Prefix}, key...)
  k := e.keyfunc(key)
	v, ok := os.LookupEnv(k)
  return v, ok
}
