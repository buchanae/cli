package main

import (
  "fmt"
  "flag"
  "github.com/buchanae/roger/roger"
  "github.com/buchanae/roger/example"
  "os"
)

import (
  "net/http/httptest"
  "io/ioutil"
  "net"
  "net/http"
)


func main() {
  c := example.DefaultConfig()
  //ts := testServer("gcemeta.json")
  //defer ts.Close()

  var configPath string
  fp := roger.Flags(c)
  fp.Flags.Init("example", flag.ExitOnError)
  fp.Flags.StringVar(&configPath, "config", configPath, "Path to config file.")
  fp.Flags.Parse(os.Args[1:])

  //gce.URL = "http://localhost:20002"
  errs := roger.Load(c,
    roger.Files(
      "/etc/roger/example.conf.yml",
      "$HOME/.example.conf.yml",
      ".example.conf.yml",
      configPath,
    ),
    roger.GCE(),
    roger.Openstack(),
    roger.Consul("127.0.0.1:8500"),
    roger.Env("example"),
    fp,
  )
  for _, e := range errs {
    fmt.Fprintln(os.Stderr, e)
  }

  c.Worker.TaskReaders.Dynamo = c.Dynamo
  c.Worker.EventWriters.Dynamo = c.Dynamo
  c.Scheduler.Worker = c.Worker
  c.Worker.TaskReaders.Dynamo = c.Dynamo
  c.Worker.EventWriters.Dynamo = c.Dynamo
  c.Worker.Storage = c.Storage

  verrs := roger.Validate(map[string]roger.Validator{
    "Dynamo": c.Dynamo,
  })
  for _, e := range verrs {
    fmt.Fprintln(os.Stderr, e)
 }

  y := roger.ToYAML(c, roger.ExcludeDefaults(example.DefaultConfig()))
  fmt.Println(y)
}

func testServer(filename string) *httptest.Server {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Start test server
	lis, err := net.Listen("tcp", ":20002")
	if err != nil {
		panic(err)
	}
	// Set up test server response
	mux := http.NewServeMux()
	mux.HandleFunc("/computeMetadata/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})
	ts := httptest.NewUnstartedServer(mux)
	ts.Listener = lis
	ts.Start()
  return ts
}
