package roger

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "testing"
  . "github.com/buchanae/roger/roger"
  "github.com/buchanae/roger/example"
)

func TestEnvKey(t *testing.T) {
  m := map[string]string{
    "one_two_three_four": "One.Two.ThreeFour",
  }
  for expected, in := range m {
    got := EnvKey(in)
    if got != expected {
      t.Errorf("expected %s, got %s", expected, got)
    }
  }
}

func TestPrefixEnvKey(t *testing.T) {
  got := PrefixEnvKey("RogerThat")("One.TwoThree")
  if got != "roger_that_one_two_three" {
    t.Errorf("expected roger_that_one_two, got %s", got)
  }
}


func TestFull(t *testing.T) {
  c := example.DefaultConfig()

  f := NewFileProvider("../example/default-config.yaml")

  os.Setenv("example_server_host_name", "set-by-env-alias")

  fs := flag.NewFlagSet("roger-example", flag.ExitOnError)
  AddFlags(c, fs, FlagKey)
  fs.Parse([]string{
    "-w", "set-by-flag-alias",
  })

  errs := Load(c,
    f,
    NewEnvProvider("example"),
    NewFlagProvider(fs),
  )

  for _, e := range errs {
    // Example of accessing name of unknown field.
    /*
    if f, ok := IsUnknownField(e); ok {
      fmt.Println(f)
    }
    */
    fmt.Println(e)
  }

  c.Scheduler.Worker = c.Worker
  c.Worker.TaskReaders.Dynamo = c.Dynamo
  c.Worker.EventWriters.Dynamo = c.Dynamo

  for _, err := range Validate(c) {
    fmt.Println(err)
 }

  if c.Server.HostName != "set-by-env-alias" {
    t.Error("expected Server.HostName to be set by env alias")
  }

  if c.Worker.WorkDir != "set-by-flag-alias" {
    t.Error("expected Worker.WorkDir to be set by flag alias")
  }

  if c.Dynamo.TableBasename != "set-by-yaml" {
    t.Error("expected Dynamo.TableBasename to be set by yaml")
  }

  if c.Worker.TaskReaders.Dynamo.TableBasename != "set-by-yaml" {
    t.Error("expected Worker.TaskReaders.Dynamo.TableBasename to be set by yaml")
  }

  if c.Worker.EventWriters.Dynamo.TableBasename != "set-by-yaml" {
    t.Error("expected Worker.EventWriters.Dynamo.TableBasename to be set by yaml")
  }

  y := ToYAML(c, example.DefaultConfig())
  s := strings.TrimSpace(y)
  if s != expectedYAML {
    t.Errorf("unexpected yaml:\n%s", s)
    t.Logf("Expected yaml:\n%s", expectedYAML)
  }
}

var expectedYAML = strings.TrimSpace(`
Server:
  # Server host name
  HostName: set-by-env-alias
  MaxExecutorLogSize: 20KB
Worker:
  # Directory to write task files to
  WorkDir: set-by-flag-alias
  # How often the worker sends task log updates
  UpdateRate: 10s
  # Max bytes to store in-memory between updates
  BufferSize: 11KB
Scheduler:
  Node:
    # How often the node sends update requests to the server.
    UpdateRate: 10s
    # RPC address of the Funnel server
    ServerAddress: localhost:9090
Dynamo:
  TableBasename: set-by-yaml
Storage:
  Local:
    AllowedDirs: [./ anotherdir]
`)
