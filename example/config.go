package example

import (
  "github.com/buchanae/roger/example/server"
  "github.com/buchanae/roger/example/worker"
  "github.com/buchanae/roger/example/scheduler"
  "github.com/buchanae/roger/example/storage"
  "github.com/buchanae/roger/example/logger"
  "github.com/buchanae/roger/example/dynamo"
)

func DefaultConfig() *Config {
  return &Config{
    Server: server.DefaultConfig(),
    Worker: worker.DefaultConfig(),
    Scheduler: scheduler.DefaultConfig(),
    Log: logger.DefaultConfig(),
    Dynamo: dynamo.DefaultConfig(),
    Storage: storage.DefaultConfig(),
  }
}

type Config struct {
  Server server.Config
  Worker worker.Config
  Scheduler scheduler.Config
  Log logger.Config
  Dynamo dynamo.Config
  Storage storage.Config
}
