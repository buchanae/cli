package main

import "time"

type Config struct {
  Server ServerConfig
  // Client configuration.
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
