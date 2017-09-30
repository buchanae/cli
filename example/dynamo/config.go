package dynamo

func DefaultConfig() Config {
  return Config{
    TableBasename: "funnel",
  }
}

// DynamoDB describes the configuration for Amazon DynamoDB backed processes
// such as the event writer and server.
type Config struct {
	Region        string
	Key           string
	Secret        string
	TableBasename string
}
