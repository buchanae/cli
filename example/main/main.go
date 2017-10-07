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
  l := roger.Loader{
    Ignore: []string{
      "Scheduler.Worker",
    },
    Alias: map[string]string{
      "host": "Server.HostName",
      "w": "Worker.WorkDir",
    },
    Files: []string{
      // TODO this assumes that this script is running
      //      from the root of the roger repo.
      "example/default-config.yaml",
    },
    EnvKeyfunc: roger.PrefixEnvKey("example"),
    FlagSet: flag.NewFlagSet("roger-example", flag.ExitOnError),
  }

  for _, e := range l.Load(c) {
    fmt.Println(e)
  }

  c.Scheduler.Worker = c.Worker
  c.Worker.TaskReaders.Dynamo = c.Dynamo
  c.Worker.EventWriters.Dynamo = c.Dynamo

  fmt.Println("worker.work_dir", c.Worker.WorkDir)

  roger.ToYAML(os.Stdout, c, l.Ignore, example.DefaultConfig())
}
