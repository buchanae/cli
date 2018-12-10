package cli

import (
	"unicode"
)

// SplitIdent splits a Go identifier, such as a function name,
// into multiple parts based on capialization.
func SplitIdent(s string) []string {
	var parts []string

	rs := []rune(s)
	start := 0

	for i := 1; i < len(s); i++ {

		prev := unicode.IsUpper(rs[i-1])
		cur := unicode.IsUpper(rs[i])

		if prev != cur {

			j := i
			if prev {
				j = i - 1
			}

			if j > start {
				parts = append(parts, string(rs[start:j]))
				start = j
			}
		}
	}
	parts = append(parts, string(rs[start:]))

	return parts
}
