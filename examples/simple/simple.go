package main

import "time"

type Opt struct {
	Server ServerOpt
	// Client configuration.
	Client ClientOpt
}

type ServerOpt struct {
	// Server name, for metadata endpoints.
	Name string
	// Address to listen on (e.g. :8080).
	Addr string
}

type ClientOpt struct {
	// Server addresses to connect to.
	Servers []string
	// Request timeout duration.
	Timeout time.Duration
}

func DefaultOpt() *Config {
	return &Opt{
		Server: ServerOpt{
			Name: "example",
			Addr: ":8080",
		},
		Client: ClientOpt{
			Servers: []string{"127.0.0.1:8080"},
			Timeout: time.Second * 30,
		},
	}
}
