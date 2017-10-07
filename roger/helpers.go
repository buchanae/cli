package roger

import (
  "fmt"
  "strings"
  "reflect"
)

type RogerVals interface {
  RogerVals() Vals
}

type Vals map[string]Val

type Val struct {
  Doc string
  val interface{}
}

func NewVal(doc string, v interface{}) Val {
  return Val{doc, v}
}

type Keyfunc func(string) string

func IdentityKey(k string) string {
  return k
}

type UnknownField string
func (u UnknownField) Error() string {
  return fmt.Sprintf("unknown field: %s", string(u))
}

func IsUnknownField(err error) (string, bool) {
  if f, ok := err.(UnknownField); ok {
    return string(f), true
  }
  return "", false
}

func FromMap(vals Vals, m map[string]interface{}) []error {
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

    if err := CoerceSet(rv, fv); err != nil {
      fmt.Println(err)
    }
  }
  return errs
}

func flatten(in map[string]interface{}, prefix string, out map[string]interface{}) {
  for k, v := range in {
    path := k
    if prefix != "" {
      path = prefix + "." + k
    }

    switch x := v.(type) {
    case map[string]interface{}:
      flatten(x, path, out)
    default:
      out[path] = v
    }
  }
}

func join(path []string, delim string, prefix string, transform func(string) string) string {
  var p []string
  if prefix != "" {
    p = append(p, prefix)
  }
  for _, i := range path {
    p = append(p, transform(i))
  }
  return strings.Join(p, delim)
}

func pathname(root reflect.Type, path []int) string {
  var name []string
  for i := 0; i < len(path); i++ {
    name = append(name, root.FieldByIndex(path[:i+1]).Name)
  }
  return strings.Join(name, ".")
}

func newpathI(base []int, add ...int) []int {
  path := append([]int{}, base...)
  return append(path, add...)
}

func tryKeyfunc(key string, kf Keyfunc, d Keyfunc) string {
  if kf == nil {
    return d(key)
  }
  return kf(key)
}
