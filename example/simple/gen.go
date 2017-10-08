package main

import "github.com/buchanae/roger/roger"

func (c *Config) RogerVals() map[string]roger.Val {
	m := map[string]roger.Val{
		"Server.Name":    roger.NewVal("Server name, for metadata endpoints.", &c.Server.Name),
		"Server.Addr":    roger.NewVal("Address to listen on (e.g.", &c.Server.Addr),
		"Client.Servers": roger.NewVal("Server addresses to connect to.", &c.Client.Servers),
		"Client.Timeout": roger.NewVal("Request timeout duration.", &c.Client.Timeout),
	}

	return m
}

