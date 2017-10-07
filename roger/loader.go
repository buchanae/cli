package roger

import (
  "fmt"
)

type Provider interface {
  Init() error
  Lookup(key string) (interface{}, error)
}

func Load(rv RogerVals, ps ...Provider) []error {
  var errs []error

  vals := rv.RogerVals()

  for _, p := range ps {
    // Just in case, avoid nil panic
    if p == nil {
      continue
    }

    if err := p.Init(); err != nil {
      errs = append(errs, err)
      continue
    }

    for k, v := range vals {
      x, err := p.Lookup(k)
      if x != nil && err == nil {
        err = CoerceSet(v, x)
      }
      if err != nil {
        errs = append(errs, fmt.Errorf("error loading %s: %s", k, err))
      }
    }
  }

  errs = append(errs, Validate(rv)...)
  return errs
}
