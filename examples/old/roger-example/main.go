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
  var gce, openstack, consul, etcd bool
  var includeDefaults, includeEmpty bool

  fp := roger.Flags(c, "example", flag.ExitOnError)
  fp.StringVar(&configPath, "config", configPath, "Path to config file.")
  fp.BoolVar(&gce, "gce-config", gce, "Enable GCE VM instance metadata config loading.")
  fp.BoolVar(&openstack, "openstack-config", openstack, "Enable Openstack VM instance metadata config loading.")
  // TODO these should take server url strings too?
  fp.BoolVar(&consul, "consul-config", consul, "Enable Consul config loading.")
  fp.BoolVar(&etcd, "etcd-config", etcd, "Enable etcd config loading.")
  fp.BoolVar(&includeDefaults, "include-defaults", includeDefaults, "Include default values in output.")
  fp.BoolVar(&includeEmpty, "include-empty", includeEmpty, "Include empty values in output.")
  fp.Parse(os.Args[1:])

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

  verrs := roger.Validate(map[string]roger.Validator{
    "Dynamo": c.Dynamo,
  })
  for _, e := range verrs {
    fmt.Fprintln(os.Stderr, e)
 }

  mar := roger.YAMLMarshaler{
    IncludeEmpty: includeEmpty,
  }
  if !includeDefaults {
    mar.ExcludeDefaults = example.DefaultConfig()
  }

  fmt.Println(mar.Marshal(c))
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
