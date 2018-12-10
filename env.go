package cli

import (
	"os"
	"strings"
)

// EnvProvider provides access to values from environment variables.
type EnvProvider struct {
	KeyFunc
	Prefix string
}

// Env returns a Provider that loads values from environment variables.
// Keys are prefixed with the given prefix.
func Env(prefix string) *EnvProvider {
	return &EnvProvider{Prefix: prefix}
}

// Init initializes the provider.
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

// Lookup returns the value of the given key.
func (e *EnvProvider) Lookup(key []string) (interface{}, bool) {
	key = append([]string{e.Prefix}, key...)
	k := e.keyfunc(key)
	v, ok := os.LookupEnv(k)
	return v, ok
}
