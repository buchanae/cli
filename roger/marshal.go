package roger

import (
  "fmt"
  "reflect"
  "strings"
  "io"
)

func ToYAML(w io.Writer, i interface{}, vals Vals, ignore []string, d interface{}) string {
  y := yamler{
    rootType: reflect.TypeOf(i).Elem(),
    rootVal: reflect.ValueOf(i).Elem(),
    defaultType: reflect.TypeOf(d).Elem(),
    defaultVal: reflect.ValueOf(d).Elem(),
    ignore: map[string]struct{}{},
    vals: vals,
    writer: w,
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
  defaultType reflect.Type
  defaultVal reflect.Value
  ignore map[string]struct{}
  includeEmpty bool
  includeDefault bool
  vals Vals
  writer io.Writer
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
    if !y.includeEmpty {
      zero := reflect.Zero(ft.Type)
      eq := reflect.DeepEqual(zero.Interface(), fv.Interface())
      //fmt.Println("EQ", name, eq)
      if eq {
        continue
      }
    }

    // Ignore default values if includeDefault is false.
    if !y.includeDefault {
      dfv := y.defaultVal.FieldByIndex(path)
      eq := reflect.DeepEqual(dfv.Interface(), fv.Interface())
      //fmt.Println("DEFAULT", name, eq)
      if eq {
        continue
      }
    }

    switch fv.Kind() {
    case reflect.Struct:
      fmt.Fprintf(y.writer, "%s%s:\n", indent, ft.Name)
      y.marshal(path)

    default:
      if v, ok := y.vals[name]; ok && v.Doc != "" {
        fmt.Fprintf(y.writer, "%s# %s\n", indent, v.Doc)
      }
      fmt.Fprintf(y.writer, "%s%s: %v\n", indent, ft.Name, fv)
    }
  }
}
