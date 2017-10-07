package example

import "github.com/buchanae/roger/roger"

func (c *Config) RogerVals() roger.Vals {
	m := map[string]roger.Val{
		"Server.Name":                        roger.NewVal("Server name.", &c.Server.Name),
		"Server.HostName":                    roger.NewVal("Server host name", &c.Server.HostName),
		"Server.HTTPPort":                    roger.NewVal("", &c.Server.HTTPPort),
		"Server.RPCPort":                     roger.NewVal("Port to serve gRPC traffic on", &c.Server.RPCPort),
		"Server.Password":                    roger.NewVal("", &c.Server.Password),
		"Server.DisableHTTPCache":            roger.NewVal("Disable http", &c.Server.DisableHTTPCache),
		"Server.MaxExecutorLogSize":          roger.NewVal("", &c.Server.MaxExecutorLogSize),
		"Server.Logger.Level":                roger.NewVal("Log level docs", &c.Server.Logger.Level),
		"Server.Logger.Formatter":            roger.NewVal("", &c.Server.Logger.Formatter),
		"Server.Logger.OutputFile":           roger.NewVal("", &c.Server.Logger.OutputFile),
		"Server.Logger.Foo.FooField":         roger.NewVal("", &c.Server.Logger.Foo.FooField),
		"Server.Logger.Foo.Level":            roger.NewVal("Foo level docs", &c.Server.Logger.Foo.Level),
		"Worker.WorkDir":                     roger.NewVal("Directory to write task files to", &c.Worker.WorkDir),
		"Worker.UpdateRate":                  roger.NewVal("How often the worker sends task log updates", &c.Worker.UpdateRate),
		"Worker.BufferSize":                  roger.NewVal("Max bytes to store in-memory between updates", &c.Worker.BufferSize),
		"Worker.Logger.Level":                roger.NewVal("Log level docs", &c.Worker.Logger.Level),
		"Worker.Logger.Formatter":            roger.NewVal("", &c.Worker.Logger.Formatter),
		"Worker.Logger.OutputFile":           roger.NewVal("", &c.Worker.Logger.OutputFile),
		"Worker.Logger.Foo.FooField":         roger.NewVal("", &c.Worker.Logger.Foo.FooField),
		"Worker.Logger.Foo.Level":            roger.NewVal("Foo level docs", &c.Worker.Logger.Foo.Level),
		"Worker.TaskReader":                  roger.NewVal("", &c.Worker.TaskReader),
		"Worker.ActiveEventWriters":          roger.NewVal("", &c.Worker.ActiveEventWriters),
		"Scheduler.ScheduleRate":             roger.NewVal("How often to run a scheduler iteration.", &c.Scheduler.ScheduleRate),
		"Scheduler.ScheduleChunk":            roger.NewVal("How many tasks to schedule in one iteration.", &c.Scheduler.ScheduleChunk),
		"Scheduler.NodePingTimeout":          roger.NewVal("How long to wait for a node ping before marking it as dead", &c.Scheduler.NodePingTimeout),
		"Scheduler.NodeInitTimeout":          roger.NewVal("How long to wait for node initialization before marking it dead", &c.Scheduler.NodeInitTimeout),
		"Scheduler.NodeDeadTimeout":          roger.NewVal("How long to wait before deleting a dead node from the DB.", &c.Scheduler.NodeDeadTimeout),
		"Scheduler.Node.ID":                  roger.NewVal("", &c.Scheduler.Node.ID),
		"Scheduler.Node.Resources.Cpus":      roger.NewVal("A Node will automatically try to detect what resources are available to it.", &c.Scheduler.Node.Resources.Cpus),
		"Scheduler.Node.Resources.RamGb":     roger.NewVal("A Node will automatically try to detect what resources are available to it.", &c.Scheduler.Node.Resources.RamGb),
		"Scheduler.Node.Resources.DiskGb":    roger.NewVal("A Node will automatically try to detect what resources are available to it.", &c.Scheduler.Node.Resources.DiskGb),
		"Scheduler.Node.Timeout":             roger.NewVal("If the node has been idle for longer than the timeout, it will shut down.", &c.Scheduler.Node.Timeout),
		"Scheduler.Node.UpdateRate":          roger.NewVal("How often the node sends update requests to the server.", &c.Scheduler.Node.UpdateRate),
		"Scheduler.Node.UpdateTimeout":       roger.NewVal("Timeout duration for UpdateNode() gRPC calls", &c.Scheduler.Node.UpdateTimeout),
		"Scheduler.Node.Metadata":            roger.NewVal("", &c.Scheduler.Node.Metadata),
		"Scheduler.Node.ServerAddress":       roger.NewVal("RPC address of the Funnel server", &c.Scheduler.Node.ServerAddress),
		"Scheduler.Node.ServerPassword":      roger.NewVal("Password for basic auth.", &c.Scheduler.Node.ServerPassword),
		"Scheduler.Node.Logger.Level":        roger.NewVal("Log level docs", &c.Scheduler.Node.Logger.Level),
		"Scheduler.Node.Logger.Formatter":    roger.NewVal("", &c.Scheduler.Node.Logger.Formatter),
		"Scheduler.Node.Logger.OutputFile":   roger.NewVal("", &c.Scheduler.Node.Logger.OutputFile),
		"Scheduler.Node.Logger.Foo.FooField": roger.NewVal("", &c.Scheduler.Node.Logger.Foo.FooField),
		"Scheduler.Node.Logger.Foo.Level":    roger.NewVal("Foo level docs", &c.Scheduler.Node.Logger.Foo.Level),
		"Scheduler.Logger.Level":             roger.NewVal("Log level docs", &c.Scheduler.Logger.Level),
		"Scheduler.Logger.Formatter":         roger.NewVal("", &c.Scheduler.Logger.Formatter),
		"Scheduler.Logger.OutputFile":        roger.NewVal("", &c.Scheduler.Logger.OutputFile),
		"Scheduler.Logger.Foo.FooField":      roger.NewVal("", &c.Scheduler.Logger.Foo.FooField),
		"Scheduler.Logger.Foo.Level":         roger.NewVal("Foo level docs", &c.Scheduler.Logger.Foo.Level),
		"Scheduler.Backend":                  roger.NewVal("", &c.Scheduler.Backend),
		"Log.Level":                          roger.NewVal("Log level docs", &c.Log.Level),
		"Log.Formatter":                      roger.NewVal("", &c.Log.Formatter),
		"Log.OutputFile":                     roger.NewVal("", &c.Log.OutputFile),
		"Log.Foo.FooField":                   roger.NewVal("", &c.Log.Foo.FooField),
		"Log.Foo.Level":                      roger.NewVal("Foo level docs", &c.Log.Foo.Level),
		"Dynamo.Region":                      roger.NewVal("", &c.Dynamo.Region),
		"Dynamo.Key":                         roger.NewVal("", &c.Dynamo.Key),
		"Dynamo.Secret":                      roger.NewVal("", &c.Dynamo.Secret),
		"Dynamo.TableBasename":               roger.NewVal("", &c.Dynamo.TableBasename),
	}
	if v, ok := m["Server.HostName"]; ok {
		m["host"] = v
	}
	if v, ok := m["Worker.WorkDir"]; ok {
		m["w"] = v
	}

	return m
}

