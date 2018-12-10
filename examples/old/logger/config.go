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

type Foo struct {
  FooField string
  // Foo level docs
  Level string
}

// Config provides configuration for a logger.
type Config struct {
  // Log level docs
	Level      string
	Formatter  string
	OutputFile string
  Foo
}

func (c Config) Validate() (errs []error) {
  if c.Level != "info" && c.Level != "error" && c.Level != "debug" {
    errs = append(errs, fmt.Errorf("unknown level %s", c.Level))
  }
  return
}
