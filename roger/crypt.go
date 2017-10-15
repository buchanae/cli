package roger

import (
  "github.com/xordataexchange/crypt/config"
)

type ConsulProvider struct {
  Servers []string
  cm config.ConfigManager
}

func NewConsulProvider(servers []string) *ConsulProvider {
  return &ConsulProvider{Servers: servers}
}

func (c *ConsulProvider) Init() error {
  cm, err := config.NewStandardConsulConfigManager(c.Servers)
  if err != nil {
    return err
  }
  c.cm = cm
  return nil
}

func (c *ConsulProvider) Lookup(key string) (interface{}, error) {
  v, err := c.cm.Get(key)
	if err != nil {
		return nil, nil
	}
	return v, nil
}
