package roger

import (
  "fmt"
)

// Validator defines the interface of a value that should be validated.
type Validator interface {
  Validate() []error
}

// ValidationError defines information about a validation error.
type ValidationError struct {
  Key string
  Err error
}

// Error returns an error message.
func (v *ValidationError) Error() string {
  return fmt.Sprintf("validation error %s: %s", v.Key, v.Err)
}

// Validate looks for values in the given map which implement Validator,
// if calls Validate, collecting ValidationErrors.
//
// TODO this will probably go away. Haven't figured out a reasonable, generic
//      interface to handling the wide range of validation cases. Probably
//      better implemented as custom, explicit logic.
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
