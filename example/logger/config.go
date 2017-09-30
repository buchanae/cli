package logger

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
