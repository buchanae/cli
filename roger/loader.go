package roger

import (
  "flag"
  "os"
)

type Loader struct {
  Ignore []string
  Alias map[string]string
  Files []string
  EnvKeyfunc Keyfunc
  FlagSet *flag.FlagSet
}

func (l *Loader) Load(rv RogerVals) []error {
  return l.LoadArgs(rv, os.Args[1:])
}

func (l *Loader) LoadArgs(rv RogerVals, args []string) []error {
  var errs []error

  vals := rv.RogerVals()
  vals.DeletePrefix(l.Ignore...)
  vals.Alias(l.Alias)

  for _, path := range l.Files {
    e := FromFile(vals, path)
    errs = append(errs, e...)
  }
  FromAllEnv(vals, l.EnvKeyfunc)

  if l.FlagSet != nil {
    AddFlags(vals, l.FlagSet)
    err := l.FlagSet.Parse(args)
    if err != nil {
      errs = append(errs, err)
    }
  }

  errs = append(errs, Validate(rv, l.Ignore)...)
  return errs
}
