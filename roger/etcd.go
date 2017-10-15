package roger

import (
  "github.com/xordataexchange/crypt/config"
)

type EtcdProvider struct {
  Servers []string
  cm config.ConfigManager
}

func Etcd(servers ...string) *EtcdProvider {
  return &EtcdProvider{Servers: servers}
}

func (c *EtcdProvider) Init() error {
  cm, err := config.NewStandardEtcdConfigManager(c.Servers)
  if err != nil {
    return err
  }
  c.cm = cm
  return nil
}

func (c *EtcdProvider) Lookup(key string) (interface{}, error) {
  v, err := c.cm.Get(key)
	if err != nil {
		return nil, nil
	}
	return v, nil
}
