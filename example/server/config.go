package server

import (
  "github.com/buchanae/roger/example/logger"
  "github.com/alecthomas/units"
)

func DefaultConfig() Config {
  return Config{
    Name: "Funnel",
    HostName: "localhost",
    HTTPPort: "8000",
    RPCPort: "9090",
    MaxExecutorLogSize: 10 * units.KB,
    DisableHTTPCache: true,
    Logger: logger.DefaultConfig(),
  }
}

type Config struct {
  Name string
	HostName    string
	HTTPPort    string
	RPCPort     string
	Password    string
	DisableHTTPCache   bool
	MaxExecutorLogSize units.MetricBytes
	Logger             logger.Config
}

// HTTPAddress returns the HTTP address based on HostName and HTTPPort
func (c Config) HTTPAddress() string {
	if c.HostName != "" && c.HTTPPort != "" {
		return "http://" + c.HostName + ":" + c.HTTPPort
	}
	return ""
}

// RPCAddress returns the RPC address based on HostName and RPCPort
func (c Config) RPCAddress() string {
	if c.HostName != "" && c.RPCPort != "" {
		return c.HostName + ":" + c.RPCPort
	}
	return ""
}
