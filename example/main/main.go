package main

import (
  "fmt"
  "flag"
  "github.com/buchanae/roger/roger"
  "github.com/buchanae/roger/example"
  "os"
)

func main() {
  c := example.DefaultConfig()

  fs := flag.NewFlagSet("example", flag.ExitOnError)
  // TODO add config file flag
  roger.AddFlags(c, fs, roger.FlagKey)

  f, _ := roger.NewFileProvider("example/default-config.yaml")
  env := &roger.EnvProvider{Prefix: "example"}

  for _, e := range l.Load(c) {
    fmt.Println(e)
  }

  c.Scheduler.Worker = c.Worker
  c.Worker.TaskReaders.Dynamo = c.Dynamo
  c.Worker.EventWriters.Dynamo = c.Dynamo

  fmt.Println("worker.work_dir", c.Worker.WorkDir)

  roger.ToYAML(os.Stdout, c, l.Ignore, example.DefaultConfig())
}
