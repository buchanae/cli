Roger is a library and code generation tool for managing access
to application configuration via flags, files, environment variables, and more.

### Alpha quality

Roger is new and experimental and rough around the edges. Handle with care.  
There are likely bugs and panics and unhandled edge cases.

### Usage

Let's say we have a configuration struct:

```go
// config.go
package main

import "time"

type Config struct {
  Server ServerConfig
  Client ClientConfig
}

type ServerConfig struct {
  // Server name, for metadata endpoints.
  Name string
  // Address to listen on (e.g. :8080).
  Addr string
}

type ClientConfig struct {
  // Server addresses to connect to.
  Servers []string
  // Request timeout duration.
  Timeout time.Duration
}

func DefaultConfig() *Config {
  return &Config{
    Server: ServerConfig{
      Name: "example",
      Addr: ":8080",
    },
    Client: ClientConfig{
      Servers: []string{"127.0.0.1:8080"},
      Timeout: time.Second * 30,
    },
  }
}
```

Get the roger command line tool:
```
go get github.com/buchanae/roger
```

Generate some data for the roger library to use:
```
roger -root Config -out gen.go ./
```

Now, write your main func:
```go
// main.go
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
  y := roger.ToYAML(c, roger.ExcludeDefaults(DefaultConfig()))
  fmt.Println(y)
}
```

Check the CLI help:
```
go run *.go -h
```

Set some values:
```
ex_server_name=set-by-env go run *.go -client.timeout 1m
```

Precedence plays a role, e.g. flags override env. vars:
```
ex_server_name=set-by-env go run *.go -server.name set-by-flag
```

Slices need single quotes and spaces:
```
go run *.go -client.servers 'srv1 srv2'
```

This example code is in [./example/simple](./example/simple).
These docs are a work in progress. There's a more complex example in [./example](./example).

### Why?

While developing [Funnel][1], we've frequently run into issues with configuration:  

- Fields which are not propery named, cased, and indented in a YAML config file
  are not reported as errors.

- Config is constantly changing, and becomes inconsistent with with CLI flags,
  YAML loading, and documentation.

- Sub-sections of the config share the same datastructure. When config is loaded
  from flags and files, careful merging is required to make sure these sections
  stay correctly in sync. This led to complex code and many bugs.

- Config was inflexible. Support for environment variables and other sources
  was not provided.

- Only a small subset of the config was available via the CLI, which is painful
  when debugging and deploying in a distributed environment on multiple cloud providers.

- Dealing with time.Duration in YAML was a pain (quick, write out 1 second in nanoseconds).

Sounds like we needed something like [Viper][2] or one of the [many other][3] config/flag
libraries out there. After surveying these, I had a few things I wanted that I thought
were missing:

- Viper looks like it will become deeply ingrained in your codebase. Outside of the
  bootstrap code, I'd like config to just be structs.

- Other libraries go overboard with struct tags, in my opinion. Again, outside of the
  bootstrap code, the effect on the rest of the codebase should be minimal.

- Since config is defined and used as structs, the bootstrap code should use those
  structs as much as possible to load data and generate flags, docs, etc.

- Defaults also come from (instances of) those structs.

- Docs come from Go code comments in the structs.

With these ideas, I started exploring via building roger. I won't claim this a good idea,
it's an experiment.

### TODO

- possibly rethink all the interfaces: provider, etc.
- better handling (ignoring) of errors during static analysis loading.
- probably remove validation.
- handle `interface{}`, `reflect.Type`, and `reflect.Value` with more care.
- properly marshal yaml/json slices/maps/etc.
- manage editing config file
- GCE metadata provider
- etcd provider
- consul provider
- dump json, env, flags
- handle map[string]string via "key=value" flag value
- explore "storage.local.allowed_dirs.append"
- pull fieldname from json tag
- ignore/alias fields via struct tag
- recognize misspelled env var
- case sensitivity
- reflect-based inspect (i.e. works without code-gen, but sacrifices docs)

Complex:
- reloading
- multiple config files with merging

Questions:
- how to handle pointers? cycles?
- how are slices handled in env vars?
- how are slices of structs handled in flags?
- how to handle unknown type wrappers, e.g. type Foo int64


### Examples

Unset by flag
```
go run example/main/main.go -config example/default-config.yaml -dynamo.table_basename=
```

[1]: https://github.com/ohsu-comp-bio/funnel
[2]: https://github.com/spf13/viper
[3]: https://github.com/ohsu-comp-bio/funnel/issues/252
