package roger

import (
  "os"
  "strings"
)

type EnvProvider struct {
  Keyfunc
}

func NewEnvProvider(prefix string) *EnvProvider {
  return &EnvProvider{Keyfunc: PrefixEnvKey(prefix)}
}

func (e *EnvProvider) Lookup(key string) (interface{}, bool) {
  key = tryKeyfunc(key, e.Keyfunc, EnvKey)
  return os.LookupEnv(key)
}

func EnvKey(k string) string {
  return join(strings.Split(k, "."), "_", "", underscore)
}

func PrefixEnvKey(prefix string) Keyfunc {
  return func(k string) string {
    if prefix != "" {
      k = prefix + "." + k
    }
    return EnvKey(k)
  }
}
