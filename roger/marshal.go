package roger

import (
  "fmt"
  "reflect"
  "strings"
)

func ToYAML(i interface{}, vals Vals, ignore []string) string {
  y := yamler{
    rootType: reflect.TypeOf(i).Elem(),
    rootVal: reflect.ValueOf(i).Elem(),
    ignore: map[string]struct{}{},
    vals: vals,
  }
  for _, i := range ignore {
    y.ignore[i] = struct{}{}
  }
  y.marshal(nil)
  return ""
}

type yamler struct {
  rootType reflect.Type
  rootVal reflect.Value
  ignore map[string]struct{}
  includeEmpty bool
  vals Vals
}

func (y *yamler) marshal(base []int) {
  t := y.rootVal.FieldByIndex(base)

  for j := 0; j < t.NumField(); j++ {
    indent := strings.Repeat("  ", len(base))
    path := newpathI(base, j)
    name := pathname(y.rootType, path)

    if _, ok := y.ignore[name]; ok {
      continue
    }

    ft := y.rootType.FieldByIndex(path)
    fv := y.rootVal.FieldByIndex(path)

    // Ignore unexported fields.
    if ft.PkgPath != "" {
      continue
    }

    // Ignore zero values if includeEmpty is false.
    zero := reflect.Zero(ft.Type)
    eq := reflect.DeepEqual(zero.Interface(), fv.Interface())
    if !y.includeEmpty && eq {
      continue
    }

    switch fv.Kind() {
    case reflect.Struct:
      fmt.Printf("%s%s:\n", indent, ft.Name)
      y.marshal(path)

    default:
      if v, ok := y.vals[name]; ok && v.Doc != "" {
        fmt.Printf("%s# %s\n", indent, v.Doc)
      }
      fmt.Printf("%s%s: %v\n", indent, ft.Name, fv)
    }
  }
}
