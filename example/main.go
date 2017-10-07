package main

import (
  "fmt"
  "flag"
  "github.com/buchanae/roger/roger"
  "os"
)

func main() {
  c := DefaultConfig()
  vals := c.RogerVals()

  // Aliases
  vals["host"] = vals["Server.HostName"]
  vals["w"] = vals["Worker.WorkDir"]

  // TODO alias helper

  // Simple, single delete
  delete(vals, "Scheduler.Worker.TaskReader")

  vals.DeletePrefix("Scheduler.Worker")
  vals.DeletePrefix("Worker.TaskReader")

  fs := flag.NewFlagSet("roger-example", flag.ExitOnError)
  roger.AddFlags(vals, fs)
  fs.PrintDefaults()
  fs.Parse(os.Args[1:])
  roger.FromAllEnvPrefix(vals, "example")

  errs := roger.FromYAMLFile(vals, "example/default-config.yaml")

  for _, e := range errs {
    if f, ok := roger.IsUnknownField(e); ok {
      fmt.Println(f)
    }
  }

  c.Scheduler.Worker = c.Worker

  fmt.Println("worker.work_dir", c.Worker.WorkDir)
  fmt.Println(c.Server.MaxExecutorLogSize)
  fmt.Println(c.Scheduler.ScheduleRate)
}
