package cli

type Consul struct {
}

func (c *Consul) Load(l *Loader) error {
  val := l.Get([]string{"consul"})
  if val == nil {
    return nil
  }
  _ = val
  return nil

/*
  conf, ok := val.(*ConsulConfig)
  if !ok {
    return fmt.Errorf(`"consul" opt must be an instance of *ConsulCOnfig, but got %T`, val)
  }

  conn, err := Connect(conf)
  if err != nil {
    return fmt.Errorf("connecting to consul: %v", err)
  }
  */
}
