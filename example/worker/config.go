package worker

import (
  "time"
  "github.com/buchanae/roger/example/logger"
  "github.com/buchanae/roger/example/dynamo"
  "github.com/buchanae/roger/example/storage"
  "github.com/alecthomas/units"
)

func DefaultConfig() Config {
  c := Config{
	  UpdateRate: time.Second * 5,
		BufferSize: 10 * units.KB,
    Storage: storage.DefaultConfig(),
    Logger: logger.DefaultConfig(),
    TaskReader: "rpc",
    ActiveEventWriters: []string{"rpc", "log"},
  }
  c.TaskReaders.Dynamo = dynamo.DefaultConfig()
  c.EventWriters.Dynamo = dynamo.DefaultConfig()
  c.EventWriters.RPC.UpdateTimeout = time.Second
  return c
}

type Config struct {
	Storage     storage.Config
	// Directory to write task files to
	WorkDir string
	// How often the worker sends task log updates
	UpdateRate time.Duration
	// Max `bytes` to store in-memory between updates
	BufferSize  units.MetricBytes
	Logger      logger.Config
	TaskReader  string
	TaskReaders struct {
		RPC struct {
			// RPC address of the Funnel server
			ServerAddress string
			// Password for basic auth. with the server APIs.
			ServerPassword string
		}
    Dynamo dynamo.Config
	}
	ActiveEventWriters []string
	EventWriters       struct {
		RPC struct {
			// RPC address of the Funnel server
			ServerAddress string
			// Password for basic auth. with the server APIs.
			ServerPassword string
			// Timeout duration for gRPC calls
			UpdateTimeout time.Duration
		}
    Dynamo dynamo.Config
	}
}
