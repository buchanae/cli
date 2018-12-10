package main

import (
  "fmt"
  "github.com/buchanae/roger/roger"
)

func main() {
  c := DefaultConfig()
  fp := roger.NewFlagProvider(c)

  roger.Load(c,
    roger.OptionalFileProvider("roger.conf.yml"),
    roger.NewEnvProvider("ex"),
    fp,
  )

  fmt.Println("servers", c.Client.Servers)

  y := roger.ToYAML(c, roger.ExcludeDefaults(DefaultConfig()))
  fmt.Println(y)
}
