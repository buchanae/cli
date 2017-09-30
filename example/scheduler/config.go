package scheduler

import (
  "time"
  "github.com/buchanae/roger/example/logger"
  "github.com/buchanae/roger/example/worker"
)

func DefaultConfig() Config {
  return Config{
    ScheduleRate:    time.Second,
    ScheduleChunk:   10,
    NodePingTimeout: time.Minute,
    NodeInitTimeout: time.Minute * 5,
    NodeDeadTimeout: time.Minute * 5,
    Node: NodeConfig{
      Timeout:       -1,
      UpdateRate:    time.Second * 5,
      UpdateTimeout: time.Second,
      Metadata:      map[string]string{},
      Logger:        logger.DefaultConfig(),
    },
    Logger: logger.DefaultConfig(),
    Worker: worker.DefaultConfig(),
  }
}

// Config contains funnel's basic scheduler configuration.
type Config struct {
	// How often to run a scheduler iteration.
	ScheduleRate time.Duration
	// How many tasks to schedule in one iteration.
	ScheduleChunk int
	// How long to wait for a node ping before marking it as dead
	NodePingTimeout time.Duration
	// How long to wait for node initialization before marking it dead
	NodeInitTimeout time.Duration
	// How long to wait before deleting a dead node from the DB.
	NodeDeadTimeout time.Duration
	// Node configuration
	Node NodeConfig
	// Logger configuration
	Logger logger.Config
  // Worker configuration
  Worker worker.Config
}

// NodeConfig contains the configuration for a node. Nodes track available resources
// for funnel's basic scheduler.
type NodeConfig struct {
	ID      string
	// A Node will automatically try to detect what resources are available to it.
	// Defining Resources in the Node configuration overrides this behavior.
	Resources struct {
		Cpus   uint32
		RamGb  float64 // nolint
		DiskGb float64
	}
	// If the node has been idle for longer than the timeout, it will shut down.
	// -1 means there is no timeout. 0 means timeout immediately after the first task.
	Timeout time.Duration
	// How often the node sends update requests to the server.
	UpdateRate time.Duration
	// Timeout duration for UpdateNode() gRPC calls
	UpdateTimeout time.Duration
	Metadata      map[string]string
	// RPC address of the Funnel server
	ServerAddress string
	// Password for basic auth. with the server APIs.
	ServerPassword string
	Logger         logger.Config
}
