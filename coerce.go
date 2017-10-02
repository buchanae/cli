package roger

import (
  "github.com/spf13/cast"
  "github.com/alecthomas/units"
  "time"
  "reflect"
)

func CoerceSet(dest interface{}, val interface{}) error {
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
    return fmt.Errorf("unknown dest type", dest)
  }

  if err != nil {
    return fmt.Errorf("error casting", l.Path, val, err)
  }

  reflect.ValueOf(ptrs[k]).Elem().Set(reflect.ValueOf(v))
  l.Value.Set(reflect.ValueOf(casted))
  return nil
}
