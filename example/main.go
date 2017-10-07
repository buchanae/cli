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

  fs := flag.NewFlagSet("roger-example", flag.ExitOnError)
  roger.AddFlags(vals, fs)
  fs.PrintDefaults()
  fs.Parse(os.Args[1:])
  roger.SetAllFromEnvPrefix(vals, "example")

  /*
  y, err := roger.LoadYAML("example/default-config.yaml")
  if err != nil {
    panic(err)
  }
  roger.SetFromMap(vals, y)
  //fmt.Println(y, err)
  */

  c.Scheduler.Worker = c.Worker

  fmt.Println("worker.work_dir", c.Worker.WorkDir)
  fmt.Println(c.Server.MaxExecutorLogSize)
  fmt.Println(c.Scheduler.ScheduleRate)
}
