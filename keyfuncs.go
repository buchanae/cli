package cli

import (
	"strings"
)

// KeyFunc defines a function which transforms a key.
type KeyFunc func([]string) string

var (
	// DotKey joins key parts with a "." and converts to lowercase.
	DotKey KeyFunc = lowerJoin(".")

	// UnderscoreKey joins key parts with a "_" and converts to lowercase.
	UnderscoreKey = lowerJoin("_")

	// DashKey joins key parts with a "-" and converts to lowercase.
	DashKey = lowerJoin("-")
)

func lowerJoin(delim string) KeyFunc {
	return func(k []string) string {
		return strings.ToLower(strings.Join(k, delim))
	}
}
