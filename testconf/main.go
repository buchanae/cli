package main

import (
  "fmt"
)

func main() {
  c := DefaultConfig()
  c.Set("Worker.WorkDir", "foo")
  fmt.Println(c.Worker.WorkDir)
}
