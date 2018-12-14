package cli

import (
	"os"
	"strings"
)

// Env loads option values from environment variables
// with the given prefix. The prefix and option keys
// are converted to uppercase.
func Env(prefix string) Provider {
	return &env{prefix}
}

// env loads option values from environment variables.
type env struct {
	Prefix string
}

func (e *env) Provide(l *Loader) error {
	for _, key := range l.Keys() {

		var prefixed []string
		if e.Prefix != "" {
			prefixed = append([]string{e.Prefix}, key...)
		} else {
			prefixed = key
		}
		k := strings.ToUpper(UnderscoreKey(prefixed))

		v, ok := os.LookupEnv(k)
		if !ok {
			continue
		}
		l.Set(key, v)
	}
	return nil
}
