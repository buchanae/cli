package cli

import (
	"os"
	"strings"
)

func Env(prefix string) Source {
  return &env{prefix}
}

// env loads option values from environment variables.
type env struct {
	Prefix string
}

func (e *env) Load(l *Loader) error {
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
