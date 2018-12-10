package cli

import (
	"strings"
)

// Keyfunc defines a function which transforms a key. For example,
//
// `EnvKey` is a Keyfunc which transforms "Root.Sub.SubOne" to "root_sub_sub_one",
// so that the key may be defined via environment variables.
//
// The input is always a key of the style "Root.Sub.SubOne";
// Go-style variable names are joined by ".".
type KeyFunc func([]string) string

func DotKey(k []string) string {
  return strings.ToLower(strings.Join(k, "."))
}

func UnderscoreKey(k []string) string {
  return strings.ToLower(strings.Join(k, "_"))
}

func DashKey(k []string) string {
  return strings.ToLower(strings.Join(k, "-"))
}

func PrefixKeyFunc(prefix string, kf KeyFunc) KeyFunc {
  return func(k []string) string {
    j := append([]string{prefix}, k...)
    return kf(j)
  }
}
