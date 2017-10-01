package logger

import (
  "fmt"
)

func DefaultConfig() Config {
  return Config{
    Level: "info",
    Formatter: "text",
  }
}

// Config provides configuration for a logger.
type Config struct {
	Level      string
	Formatter  string
	OutputFile string
}

func (c Config) Validate() (errs []error) {
  if c.Level != "info" && c.Level != "error" && c.Level != "debug" {
    errs = append(errs, fmt.Errorf("unknown level %s", c.Level))
  }
  return
}
