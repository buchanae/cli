[![GoDoc](https://godoc.org/github.com/buchanae/cli?status.svg)](https://godoc.org/github.com/buchanae/cli)

`cli` helps streamline CLI and configuration code in Go applications.

`cli` uses static analysis and code generation to make building a CLI
less tedious. Config/options are described by convention, which allows
mostly automatic generation of flags, docs, loaders, and other utilities.

`cli` is my attempt to solve the issues I repeatedly encounter while
building Go applications.

# Status

`cli` is beta quality. It's new and hasn't been fully proven yet.

# Usage

Install:
```
go get github.com/buchanae/cli/cmd/cli
```

Write a CLI in `main_cli.go` (the `_cli.go` suffix is important):
```go
package main

import (
  "fmt"
  "net/http"
  "github.com/buchanae/cli"
)

//go:generate cli .

type ServerOpt struct {
  // Server name.
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

// Run a simple echo server.
// Example: ./server run --name "my-server" "Hello, world!"
func Run(opt ServerOpt, msg string) {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "from server %q: %s\n", opt.Name, msg)
  })
  cli.Check(http.ListenAndServe(opt.Addr, nil))
}

func main() {
  cli.AutoCobra("server", specs())
}
```

Generate the CLI command handlers and build it:
```
go generate .
found cli "Run"
generated file generated_specs.go

go build -o server .
```

Try it out:
```
./server run 'Hello, world!'

curl localhost:8080
from server "cli-example": Hello, world!
```

Try it with some different options:
```
./server run --name "my-server" --addr :8081 'Hello, world!'

curl localhost:8081
from server "my-server": Hello, world!
```

Try it with a config file:
```
cat config.yaml
addr: ":8081"
name: "config-server"

./server run 'Hi!'

curl localhost:8081
from server "config-server": Hi!
```

This example code is in [./example/server](./example/server).

# Conventions

- Only files with the suffix `_cli.go` are analyzed.
- Exported functions are turned into CLI commands.
- If a function has a struct-type argument named `opt`,
  the fields of that type are used to generate flags, docs,
  file loaders, etc.
- If an option type `MyOpt` has a matching function `DefaultMyOpt() MyOpt`,
  that function will provide default values for the options.
- Command function arguments are coerced from CLI positional arguments,  
  e.g. `Age(name string, age int)` maps to `./app age "Alex" 33`

# Helpers

I'm most familiar with [cobra][cobra] and YAML config files, so I wrote
`AutoCobra(appName string, specs []Spec)` to handle my common use age,
but hopefully `cli` is flexible enough to handle a wide variety of preferences.

# Why?

Building powerful configuration and commandline interfaces is important,
yet writing the code is tedious, error-prone, and sometimes tricky.

These are some issues I frequently encounter:

- Loading and merging config files, defaults, and CLI flags
  is error-prone and tricky.

- Config key names and CLI flag names can have inconsistent
  naming, casing, patters (http.port vs --server_port).

- Config files and CLI flags can get out of sync,
  e.g. add a new config option but forget to add a CLI flag.

- Only a subset of config is available via CLI/env flags,
  leading to lots of tedious throw away config files.

- Config can be misspelled or incorrectly formatted (e.g. indentation)
  leading to subtle behavior that is difficult to debug.

- time.Duration (and friends) is not handled well by common
  marshalers, such as YAML.

- Config docs easily get out of sync with the actual structures.

- Evolving config leads to broken systems when upgrading to newer versions.

- Writing CLI and config code is often tedious, verbose, and covered in boilerplate.
  This is especially annoying when you build a CLI with lots of commands.

- Unit testing is tricky because the CLI/config code usually interacts
  with the entire application. Organizing mocks or other tricks gets messy.

- Unit testing is tricky because, again, you want to replace the usual
  stdin/out/err with something you can test.

- Every new project develops a new pattern for validation, testing, etc.

- Wrappers for time.Duration are needed for proper (un)marshling from
  config files, but the wrappers then invade the whole codebase, leading
  to code like `time.NewTimer(time.Duration(config.TimeoutDuration))`.

Other wants:

- I want to be able to write out a YAML config file based on an instance
  of a config struct type. This is useful for docs, CLI utilities, debugging,
  bug reports, and more.

# Design Goals

1. Configuration should feel natural when created and used in code.
   Config should be based on struct types, defaults should be provided
   by instances, or functions that return instances, of those types.

1. Documentation should be written in code, as with nearly all Go documentation.
   Tools should be provided to generate other forms of documentation.

1. Flags, environment variables, config files, and other types of data sources
   should use struct types as the source of truth while loading.

1. Common config errors, such as misspelling or non-existent keys, should be caught
   by the core library.

1. Merging multiple config sources should be handled simply and robustly
   by the core library.

1. Projects should be consistent in their config, CLI, and test code.
   The pattern and tools developed here should be robust enough to support
   many projects.

1. The pattern should allow for easily removing code that is duplicated
   amongst multiple commands, such as database init code.

1. The pattern should help reduce boilerplate and other sources of verbose code.

# Design Decisions

1. Use `cli.Fatal` in top-level CLI code to minimize error checks.  I totally
   agree with Go's approach to error handling by value, but panic/recover is
   useful too. A CLI command is top-level; it isn't expected to be called from
   other code, so the error values aren't being checked by anything, and panics
   aren't escaping into general code.  CLI commands aren't being called in
   loops, so panic/recover performance isn't an issue. When an error occurs in
   a CLI command, you usually want the whole program to stop with an error,
   which seems like a good fit for Fatal/panic.
   
1. Use code generation to inspect commands and config. The best way to keep
   config and docs up to date is to have it written right next to the code in
   the form of code comments.

1. Keep code generation minimal; generate just enough information for libraries
   to do the rest at runtime. Details such as cobra command building, doc
   parsing, flag building, etc. *could* all happen during code generation, but
   it feels slightly less flexible and more likely to become complex. Also,
   more strings/data being generated as code means larger binaries for projects
   with lots of commands. Honestly, I'm on the fence here though.

1. Allow an alternative to struct tags. Sometimes you don't have access to
   the struct type, or you don't want to modify it (maybe it's generated by protoc).

# To do / Known Issues

- be able to hide/ignore fields without using a struct tag,
  for fields which you don't have access to or don't want to modify
  with cli tags.
- provide `sensitive` tag for passwords and other sensitive fields.
- properly marshal yaml/json slices/maps/etc.
- GCE metadata, etcd, consul, openstack provider
- dump json, env, flags
- handle map[string]string via "key=value" flag value
- pull fieldname from json tag
- ignore/alias fields via struct tag
- recognize misspelled env var
- case sensitivity
- manage editing config file

Complex:
- reloading
- multiple config files with merging

Questions:
- how to handle pointers? cycles?
- how are slices handled in env vars?
- how are slices of structs handled in flags?
- how to handle unknown type wrappers, e.g. type Foo int64

[viper]: https://github.com/spf13/viper
