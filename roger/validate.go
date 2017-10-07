package roger

import (
  "fmt"
  "reflect"
)

type Validator interface {
  Validate() []error
}

type ValidationError struct {
  Field string
  Err error
}

func (v *ValidationError) Error() string {
  return fmt.Sprintf("validation error %s: %s", v.Field, v.Err)
}

func Validate(i RogerVals) []error {
  v := validator{
    rootType: reflect.TypeOf(i).Elem(),
    rootVal: reflect.ValueOf(i).Elem(),
    vals: i.RogerVals(),
  }
  v.validate(nil)
  return v.errors
}

type validator struct {
  rootType reflect.Type
  rootVal reflect.Value
  errors []error
  vals Vals
}

func (v *validator) validate(base []int) {
  t := v.rootVal.FieldByIndex(base)

  for j := 0; j < t.NumField(); j++ {
    path := newpathI(base, j)
    name := pathname(v.rootType, path)

    if _, ok := v.vals[name]; !ok {
      continue
    }

    ft := v.rootType.FieldByIndex(path)
    fv := v.rootVal.FieldByIndex(path)

    // Ignore unexported fields.
    if ft.PkgPath != "" {
      continue
    }

    if x, ok := fv.Interface().(Validator); ok {
      for _, err := range x.Validate() {
        v.errors = append(v.errors, &ValidationError{
          Field: name,
          Err: err,
        })
      }
    }

    if fv.Kind() == reflect.Struct {
      v.validate(path)
    }
  }
}
