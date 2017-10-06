package main

import (
  "flag"
  "github.com/buchanae/roger/roger"
)

func main() {
  c := DefaultConfig()
  fs := flag.NewFlagSet("roger-example", flag.ExitOnError)
  roger.AddFlags(fs, c)
  fs.PrintDefaults()
}
