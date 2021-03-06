package roger

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kr/text"
	"reflect"
	"strings"
)

// ToYAML is a utility that writes the given RogerVals to a YAML string.
//
// Goals:
//   - write a reasonable subset of YAML
//   - include comments for fields
//   - be able to write only the fields which are not a zero value (see IncludeEmpty)
//   - be able to write only fields which differ from the default values (see ExcludeDefaults)
//
// This is not a full-blown YAML marshaler, likely has many unsupported edge cases,
// and is likely buggy, but if used with care it can be useful.

type YAMLMarshaler struct {
  ExcludeDefaults interface{}
  IncludeEmpty bool
}

func (y *YAMLMarshaler) Marshal(rv RogerVals) string {
  yr := yamler{
		rootType: reflect.TypeOf(rv).Elem(),
		rootVal:  reflect.ValueOf(rv).Elem(),
		vals:     rv.RogerVals(),
    includeEmpty: y.IncludeEmpty,
    defaults: y.ExcludeDefaults,
  }
  return yr.marshal(nil)
}

type yamler struct {
	rootType     reflect.Type
	rootVal      reflect.Value
	includeEmpty bool
	defaults     interface{}
	vals         map[string]Val
}

func (y *yamler) marshal(base []int) string {
	t := y.rootVal.FieldByIndex(base)
	var s string

	for j := 0; j < t.NumField(); j++ {
		indent := strings.Repeat("  ", len(base))
		path := newpathI(base, j)
		name := pathname(y.rootType, path)

		ft := y.rootType.FieldByIndex(path)
		fv := y.rootVal.FieldByIndex(path)

		// Ignore unexported fields.
		if ft.PkgPath != "" {
			continue
		}

		// Ignore zero values if includeEmpty is false.
		if !y.includeEmpty {
			zero := reflect.Zero(ft.Type)
			eq := reflect.DeepEqual(zero.Interface(), fv.Interface())
			//fmt.Println("EQ", name, eq)
			if eq {
				continue
			}
		}

		// Exclude default values if defaults is set.
		if y.defaults != nil {
			dv := reflect.ValueOf(y.defaults).Elem()
			dfv := dv.FieldByIndex(path)
			eq := reflect.DeepEqual(dfv.Interface(), fv.Interface())
			//fmt.Println("DEFAULT", name, eq)
			if eq {
				continue
			}
		}

		switch fv.Kind() {
		case reflect.Struct:
			sub := y.marshal(path)
			if sub != "" {
				s += fmt.Sprintf("%s%s:\n%s", indent, ft.Name, sub)
			}

		default:
			if v, ok := y.vals[name]; ok {
				valueString := ""

				switch x := fv.Interface().(type) {
				case uint, uint32, uint64, int, int32, int64, float32, float64,
					bool, string, fmt.Stringer:
					valueString = fmt.Sprint(x)
				default:
					b, _ := yaml.Marshal(x)
					valueString = string(b)
				}

				if strings.ContainsRune(valueString, '\n') {
					valueString = text.Indent("\n"+valueString, indent)
				}

				sub := fmt.Sprintf("%s%s: %s\n", indent, ft.Name, valueString)
				if v.Doc != "" {
					sub = fmt.Sprintf("%s# %s\n%s", indent, v.Doc, sub)
				}
				s += sub
			}
		}
	}
	return s
}
