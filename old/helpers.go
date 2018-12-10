package roger

import (
	"fmt"
  "os"
	"reflect"
	"strings"
  "regexp"
)

var rx = regexp.MustCompile("[^a-zA-Z0-9]")

// Keyfunc defines a function which transforms a key. For example,
//
// `EnvKey` is a Keyfunc which transforms "Root.Sub.SubOne" to "root_sub_sub_one",
// so that the key may be defined via environment variables.
//
// The input is always a key of the style "Root.Sub.SubOne";
// Go-style variable names are joined by ".".
type Keyfunc func(string) string

// IdentityKey is a Keyfunc that returns the input string "k" without modification.
func IdentityKey(k string) string {
	return k
}

func DotKey(k string) string {
	return join(strings.Split(k, "."), ".", underscoreName)
}

func UnderscoreKey(k string) string {
	return join(strings.Split(k, "."), "_", underscoreName)
}

func NormalizeKey(k string) string {
  k = rx.ReplaceAllString(k, ".")
  k = DotKey(k)
  k = strings.ToLower(k)
  return k
}

func PrefixKey(prefix string, kf Keyfunc) Keyfunc {
  return func(k string) string {
		if prefix != "" {
			k = prefix + "." + k
		}
    return kf(k)
  }
}

// UnknownField may be returned by config loading code (e.g. from a file)
// to signal that an unknown field was found.
type UnknownField string

// Error returns an error message.
func (u UnknownField) Error() string {
	return fmt.Sprintf("unknown field: %s", string(u))
}

// IsUnknownField returns true if the given "err" is of type UnknownField.
// If true, it also returns the error message.
func IsUnknownField(err error) (string, bool) {
	if f, ok := err.(UnknownField); ok {
		return string(f), true
	}
	return "", false
}

// FromMap is currently unused and likely broken.
// TODO revive this.
func FromMap(vals map[string]Val, m map[string]interface{}) []error {
	var errs []error
	f := map[string]interface{}{}
	flatten(m, "", f)

	for fk, fv := range f {
		// TODO
		// If there's a block defined but all its values are commented out,
		// this will show up as unknown. Debatable what should be done in that case.
		// It isn't technically unknown, but it's not very clean either.
		if fv == nil {
			continue
		}

		rv, ok := vals[fk]
		if !ok {
			errs = append(errs, UnknownField(fk))
			continue
		}

		if err := coerceSet(rv.val, fv); err != nil {
			fmt.Println(err)
		}
	}
	return errs
}

// pathname returns a string key (e.g. "Root.Sub.SubOne") for a struct field
// given by an int path (see reflect.FieldByIndex).
func pathname(root reflect.Type, path []int) string {
	var name []string
	for i := 0; i < len(path); i++ {
		name = append(name, root.FieldByIndex(path[:i+1]).Name)
	}
	return strings.Join(name, ".")
}

// newpathI helps copy a struct field index and append new entries.
func newpathI(base []int, add ...int) []int {
	path := append([]int{}, base...)
	return append(path, add...)
}
