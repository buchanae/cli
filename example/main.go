package main

import (
  "fmt"
  "flag"
  "github.com/buchanae/roger/roger"
  "os"
)

func main() {
  c := DefaultConfig()
  fs := flag.NewFlagSet("roger-example", flag.ExitOnError)
  roger.AddFlags(c, fs)
  fs.PrintDefaults()
  fs.Parse(os.Args[1:])
  roger.SetAllFromEnvPrefix(c, "example")

  y, err := roger.LoadYAML("example/default-config.yaml")
  if err != nil {
    panic(err)
  }
  roger.SetFromMap(c, y)
  //fmt.Println(y, err)

  fmt.Println("worker.work_dir", c.Worker.WorkDir)
}
