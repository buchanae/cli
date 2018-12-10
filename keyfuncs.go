package cli

import (
	"strings"
)

// KeyFunc defines a function which transforms a key.
type KeyFunc func([]string) string

// DotKey joins key parts with a "." and converts to lowercase.
func DotKey(k []string) string {
	return strings.ToLower(strings.Join(k, "."))
}

// UnderscoreKey joins key parts with a "_" and converts to lowercase.
func UnderscoreKey(k []string) string {
	return strings.ToLower(strings.Join(k, "_"))
}

// DashKey joins key parts with a "-" and converts to lowercase.
func DashKey(k []string) string {
	return strings.ToLower(strings.Join(k, "-"))
}

// PrefixKeyFunc returns a new KeyFunc which will prefix the parts
// with "prefix" and then run "kf".
func PrefixKeyFunc(prefix string, kf KeyFunc) KeyFunc {
	return func(k []string) string {
		j := append([]string{prefix}, k...)
		return kf(j)
	}
}
