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
  fp := roger.NewFlagProvider(c)
  fp.Flags.Init("example", flag.ExitOnError)
  fp.Flags.StringVar(&configPath, "config", configPath, "Path to config file.")
  fp.Flags.Parse(os.Args[1:])

  errs := roger.Load(c,
    roger.OptionalFileProvider(".example.conf.yml"),
    roger.NewFileProvider(configPath),
    roger.NewEnvProvider("example"),
    fp,
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

  y := roger.ToYAML(c, roger.ExcludeDefaults(example.DefaultConfig()))
  fmt.Println(y)
}
