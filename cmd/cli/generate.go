package main

import (
	"flag"
  "log"
  "github.com/buchanae/cli/inspect"
)

func init() {
	log.SetFlags(0)
}

func main() {
  var verbose bool
	flag.BoolVar(&verbose, "v", verbose, "Verbose logging.")
	flag.Parse()

  pkg, err := inspect.Inspect(flag.Args(), verbose)
  if err != nil {
    log.Fatal(err)
  }

  err = inspect.Generate(pkg, inspect.DefaultTemplate)
  if err != nil {
    log.Fatal(err)
  }
}
