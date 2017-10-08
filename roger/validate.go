package roger

import (
  "fmt"
)

type Validator interface {
  Validate() []error
}

type ValidationError struct {
  Key string
  Err error
}

func (v *ValidationError) Error() string {
  return fmt.Sprintf("validation error %s: %s", v.Key, v.Err)
}

func Validate(vs map[string]Validator) []*ValidationError {
  var errs []*ValidationError
  for k, v := range vs {
    for _, err := range v.Validate() {
      errs = append(errs, &ValidationError{
        Key: k,
        Err: err,
      })
    }
  }
  return errs
}
