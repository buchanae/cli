package roger

import (
  "bytes"
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
  vals := c.RogerVals()

  ignore := []string{
    "Scheduler.Worker",
    "Worker.TaskReaders.Dynamo",
    "Worker.EventWriters.Dynamo",
  }

  vals.Alias(map[string]string{
    "host": "Server.HostName",
    "w": "Worker.WorkDir",
  })

  vals.DeletePrefix(ignore...)

  errs := FromFile(vals, "../example/default-config.yaml")
  for _, e := range errs {
    if f, ok := IsUnknownField(e); ok {
      // Example of accessing name of unknown field.
      fmt.Println(f)
    } else {
      fmt.Println(e)
    }
  }

  os.Setenv("example_server_host_name", "set-by-env-alias")

  FromAllEnvPrefix(vals, "example")

  fs := flag.NewFlagSet("roger-example", flag.ExitOnError)
  AddFlags(vals, fs)
  fs.Parse([]string{
    "-w", "set-by-flag-alias",
  })

  c.Scheduler.Worker = c.Worker
  c.Worker.TaskReaders.Dynamo = c.Dynamo
  c.Worker.EventWriters.Dynamo = c.Dynamo

  for _, err := range Validate(c, ignore) {
    fmt.Println(err)
  }

  if c.Server.HostName != "set-by-env-alias" {
    t.Error("expected Server.HostName to be set by env alias")
  }

  if c.Worker.WorkDir != "set-by-flag-alias" {
    t.Error("expected Worker.WorkDir to be set by flag alias")
  }

  var y bytes.Buffer
  ToYAML(&y, c, vals, ignore, example.DefaultConfig())
  s := strings.TrimSpace(y.String())
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
  Storage:
    Local:
      AllowedDirs: [./ anotherdir]
  # Directory to write task files to
  WorkDir: set-by-flag-alias
  # How often the worker sends task log updates
  UpdateRate: 10s
  # Max bytes to store in-memory between updates
  BufferSize: 11KB
  TaskReaders:
    RPC:
      # RPC address of the Funnel server
      ServerAddress: localhost:8000
  EventWriters:
    RPC:
      # RPC address of the Funnel server
      ServerAddress: localhost:8000
Scheduler:
  Node:
    # How often the node sends update requests to the server.
    UpdateRate: 10s
    # RPC address of the Funnel server
    ServerAddress: localhost:9090
`)
