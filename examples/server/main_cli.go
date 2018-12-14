// main_cli.go
package main

import (
  "fmt"
  "net/http"
  "github.com/buchanae/cli"
)

//go:generate cli .

type ServerOpt struct {
  // Server name, for metadata endpoints.
  Name string
  // Address to listen on.
  Addr string
}

func DefaultServerOpt() ServerOpt {
  return ServerOpt{
    Name: "cli-example",
    Addr: ":8080",
  }
}

func Run(opt ServerOpt, msg string) {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "from server %q: %s\n", opt.Name, msg)
  })
  cli.Check(http.ListenAndServe(opt.Addr, nil))
}

func main() {
  cli.AutoCobra("hello-world", specs())
}
