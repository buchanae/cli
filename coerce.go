package cli

import (
	"fmt"
	"github.com/spf13/cast"
	"time"
)

// Coerce attempts to coerce "val" to the type of "dest".
func Coerce(dest interface{}, val interface{}) error {

	switch z := dest.(type) {
	case *int:
		casted, err := cast.ToIntE(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *int64:
		casted, err := cast.ToInt64E(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *int32:
		casted, err := cast.ToInt32E(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *float32:
		casted, err := cast.ToFloat32E(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *float64:
		casted, err := cast.ToFloat64E(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *bool:
		casted, err := cast.ToBoolE(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *string:
		casted, err := cast.ToStringE(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *[]string:
		casted, err := cast.ToStringSliceE(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *map[string]string:
		casted, err := cast.ToStringMapStringE(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *[]int:
		casted, err := cast.ToIntSliceE(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	case *time.Duration:
		casted, err := cast.ToDurationE(val)
		if err != nil {
			return err
		}
		*z = casted
		return nil
	}
	return fmt.Errorf("cannot coerce %T to %T, unknown type %T", val, dest, dest)
}
