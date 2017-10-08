package roger

import (
  "fmt"
  "github.com/spf13/cast"
  "github.com/alecthomas/units"
  "time"
  "reflect"
)

// coerceSet attempts to coerce "val" to the type of "dest".
// If coercion succeeds, "dest" is set to the coerced value of "val".
// coerceSet panics if "dest" is not a pointer.
func coerceSet(dest interface{}, val interface{}) error {
  var casted interface{}
  var err error

  switch dest.(type) {
  case *int:
    casted, err = cast.ToIntE(val)
  case *int64:
    casted, err = cast.ToInt64E(val)
  case *int32:
    casted, err = cast.ToInt32E(val)
  case *float32:
    casted, err = cast.ToFloat32E(val)
  case *float64:
    casted, err = cast.ToFloat64E(val)
  case *bool:
    casted, err = cast.ToBoolE(val)
  case *string:
    casted, err = cast.ToStringE(val)
  case *[]string:
    casted, err = cast.ToStringSliceE(val)
  case *units.MetricBytes:
    if s, ok := val.(string); ok {
      casted, err = units.ParseMetricBytes(s)
    }
  case *time.Duration:
    casted, err = cast.ToDurationE(val)
  default:
    return fmt.Errorf("unknown dest type: %s", dest)
  }

  if err != nil {
    return fmt.Errorf("error casting: %s", err)
  }

  reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(casted))
  return nil
}
