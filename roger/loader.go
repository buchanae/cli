package roger

import (
  "fmt"
)

type Provider interface {
  Lookup(key string) (interface{}, bool)
}

func Load(rv RogerVals, ps ...Provider) []error {
  var errs []error

  vals := rv.RogerVals()

  for _, p := range ps {
    for k, v := range vals {
      x, ok := p.Lookup(k)
      if ok {
        err := CoerceSet(v, x)
        if err != nil {
          errs = append(errs, fmt.Errorf("error loading %s: %s", k, err))
        }
      }
    }
  }

  errs = append(errs, Validate(rv)...)
  return errs
}
