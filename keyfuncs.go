package cli

import (
	"strings"
)

// KeyFunc defines a function which transforms a key.
type KeyFunc func([]string) string

var (
  // DotKey joins key parts with a "." and converts to lowercase.
  DotKey KeyFunc = func(k []string) string {
	  return strings.ToLower(strings.Join(k, "."))
  }

  // UnderscoreKey joins key parts with a "_" and converts to lowercase.
  UnderscoreKey = func(k []string) string {
	  return strings.ToLower(strings.Join(k, "_"))
  }
)

// DashKeyFunc joins key parts with a "-" and converts to lowercase.
func DashKeyFunc(k []string) string {
	return strings.ToLower(strings.Join(k, "-"))
}
