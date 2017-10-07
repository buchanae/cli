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

  ignore := []string{
    "Scheduler.Worker",
  }

  vals.Alias(map[string]string{
    "host": "Server.HostName",
    "w": "Worker.WorkDir",
  })

  vals.DeletePrefix(ignore...)

  errs := roger.FromFile(vals, "example/default-config.yaml")
  for _, e := range errs {
    if f, ok := roger.IsUnknownField(e); ok {
      // Example of accessing name of unknown field.
      fmt.Println(f)
    } else {
      fmt.Println(e)
    }
  }

  roger.FromAllEnvPrefix(vals, "example")

  fs := flag.NewFlagSet("roger-example", flag.ExitOnError)
  roger.AddFlags(vals, fs)
  fs.Parse(os.Args[1:])

  c.Scheduler.Worker = c.Worker

  for _, err := range roger.Validate(c, ignore) {
    fmt.Println(err)
  }

  fmt.Println("worker.work_dir", c.Worker.WorkDir)
}
