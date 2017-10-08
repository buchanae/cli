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

  var configPath string
  fs := flag.NewFlagSet("example", flag.ExitOnError)
  fs.StringVar(&configPath, "config", configPath, "Path to config file.")
  roger.AddFlags(c, fs, roger.FlagKey)
  fs.Parse(os.Args[1:])

  errs := roger.Load(c,
    roger.OptionalFileProvider(".example.conf.yml"),
    roger.NewFileProvider(configPath),
    roger.NewEnvProvider("example"),
    roger.NewFlagProvider(fs),
  )
  for _, e := range errs {
    fmt.Fprintln(os.Stderr, e)
  }

  c.Worker.TaskReaders.Dynamo = c.Dynamo
  c.Worker.EventWriters.Dynamo = c.Dynamo
  c.Scheduler.Worker = c.Worker
  c.Worker.TaskReaders.Dynamo = c.Dynamo
  c.Worker.EventWriters.Dynamo = c.Dynamo
  c.Worker.Storage = c.Storage

  verrs := roger.Validate(map[string]roger.Validator{
    "Dynamo": c.Dynamo,
  })
  for _, e := range verrs {
    fmt.Fprintln(os.Stderr, e)
 }

  fmt.Println(roger.ToYAML(c, example.DefaultConfig()))
}
