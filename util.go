package cli

import (
	"os"
	"unicode"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// walk through a nested map, setting option values for the leaves.
func flatten2(in map[string]interface{}, l *Loader, prefix []string) {
	for k, v := range in {
		path := append(prefix[:], k)

		switch x := v.(type) {
		case map[string]interface{}:
			flatten2(x, l, path)
		default:
			l.Set(path, v)
		}
	}
}

// splitIdent splits a Go identifier, such as a function name,
// into multiple parts based on capialization.
func splitIdent(s string) []string {
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
