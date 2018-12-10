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
	flag.Parse()

  pkg, err := inspect.Inspect(flag.Args())
  if err != nil {
    log.Fatal(err)
  }

  err = inspect.Generate(pkg, inspect.DefaultTemplate)
  if err != nil {
    log.Fatal(err)
  }
}
