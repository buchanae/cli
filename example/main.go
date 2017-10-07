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
  roger.AddFlags(fs, c)
  fs.PrintDefaults()
  roger.Dump(c)
  roger.Dump(c)
  fmt.Println(fs.Parse(os.Args[1:]))
  fmt.Println(&c.Worker.WorkDir, "worker.work_dir", c.Worker.WorkDir)
}
