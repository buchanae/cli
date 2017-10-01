package dynamo

import (
  "fmt"
)

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

func (c Config) Validate() (errs []error) {
  if c.Region == "" {
    errs = append(errs, fmt.Errorf("region is empty"))
  }
  if c.Key == "" {
    errs = append(errs, fmt.Errorf("key is empty"))
  }
  if c.Secret == "" {
    errs = append(errs, fmt.Errorf("secret is empty"))
  }
  if c.TableBasename == "" {
    errs = append(errs, fmt.Errorf("table basename is empty"))
  }
  return
}
