package cli

import (
	"fmt"
	"github.com/spf13/cast"
	"time"
)

// coerceSet attempts to coerce "val" to the type of "dest".
// If coercion succeeds, "dest" is set to the coerced value of "val".
// coerceSet panics if "dest" is not a pointer.
func coerceSet(dest interface{}, val interface{}) error {

	switch z := dest.(type) {
	case *int:
		casted, err := cast.ToIntE(val)
    if err != nil {
      return err
    }
    *z = casted
	case *int64:
		casted, err := cast.ToInt64E(val)
    if err != nil {
      return err
    }
    *z = casted
	case *int32:
		casted, err := cast.ToInt32E(val)
    if err != nil {
      return err
    }
    *z = casted
	case *float32:
		casted, err := cast.ToFloat32E(val)
    if err != nil {
      return err
    }
    *z = casted
	case *float64:
		casted, err := cast.ToFloat64E(val)
    if err != nil {
      return err
    }
    *z = casted
	case *bool:
		casted, err := cast.ToBoolE(val)
    if err != nil {
      return err
    }
    *z = casted
	case *string:
		casted, err := cast.ToStringE(val)
    if err != nil {
      return err
    }
    *z = casted
	case *[]string:
		casted, err := cast.ToStringSliceE(val)
    if err != nil {
      return err
    }
    *z = casted
	case *time.Duration:
		casted, err := cast.ToDurationE(val)
    if err != nil {
      return err
    }
    *z = casted
	default:
		return fmt.Errorf("unknown dest type: %s", dest)
	}
	return nil
}
