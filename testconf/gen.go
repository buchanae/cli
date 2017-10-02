package main

import "fmt"
import "reflect"
import "github.com/buchanae/roger"

func (c *Config) ptrs() map[string]interface{} {
  return map[string]interface{}{
    "Server.Name": &c.Server.Name,
    "Server.HostName": &c.Server.HostName,
    "Server.HTTPPort": &c.Server.HTTPPort,
    "Server.RPCPort": &c.Server.RPCPort,
    "Server.Password": &c.Server.Password,
    "Server.DisableHTTPCache": &c.Server.DisableHTTPCache,
    "Server.MaxExecutorLogSize": &c.Server.MaxExecutorLogSize,
    "Server.Logger.Level": &c.Server.Logger.Level,
    "Server.Logger.Formatter": &c.Server.Logger.Formatter,
    "Server.Logger.OutputFile": &c.Server.Logger.OutputFile,
    "Server.Logger.Foo.FooField": &c.Server.Logger.Foo.FooField,
    "Server.Logger.Foo.Level": &c.Server.Logger.Foo.Level,
    "Worker.Storage.Local.AllowedDirs": &c.Worker.Storage.Local.AllowedDirs,
    "Worker.Storage.S3.Key": &c.Worker.Storage.S3.Key,
    "Worker.Storage.S3.Secret": &c.Worker.Storage.S3.Secret,
    "Worker.Storage.S3.FromEnv": &c.Worker.Storage.S3.FromEnv,
    "Worker.Storage.GS": &c.Worker.Storage.GS,
    "Worker.Storage.Swift.UserName": &c.Worker.Storage.Swift.UserName,
    "Worker.Storage.Swift.Password": &c.Worker.Storage.Swift.Password,
    "Worker.Storage.Swift.AuthURL": &c.Worker.Storage.Swift.AuthURL,
    "Worker.Storage.Swift.TenantName": &c.Worker.Storage.Swift.TenantName,
    "Worker.Storage.Swift.TenantID": &c.Worker.Storage.Swift.TenantID,
    "Worker.Storage.Swift.RegionName": &c.Worker.Storage.Swift.RegionName,
    "Worker.WorkDir": &c.Worker.WorkDir,
    "Worker.UpdateRate": &c.Worker.UpdateRate,
    "Worker.BufferSize": &c.Worker.BufferSize,
    "Worker.Logger.Level": &c.Worker.Logger.Level,
    "Worker.Logger.Formatter": &c.Worker.Logger.Formatter,
    "Worker.Logger.OutputFile": &c.Worker.Logger.OutputFile,
    "Worker.Logger.Foo.FooField": &c.Worker.Logger.Foo.FooField,
    "Worker.Logger.Foo.Level": &c.Worker.Logger.Foo.Level,
    "Worker.TaskReader": &c.Worker.TaskReader,
    "Worker.TaskReaders.RPC.ServerAddress": &c.Worker.TaskReaders.RPC.ServerAddress,
    "Worker.TaskReaders.RPC.ServerPassword": &c.Worker.TaskReaders.RPC.ServerPassword,
    "Worker.TaskReaders.Dynamo.Region": &c.Worker.TaskReaders.Dynamo.Region,
    "Worker.TaskReaders.Dynamo.Key": &c.Worker.TaskReaders.Dynamo.Key,
    "Worker.TaskReaders.Dynamo.Secret": &c.Worker.TaskReaders.Dynamo.Secret,
    "Worker.TaskReaders.Dynamo.TableBasename": &c.Worker.TaskReaders.Dynamo.TableBasename,
    "Worker.ActiveEventWriters": &c.Worker.ActiveEventWriters,
    "Worker.EventWriters.RPC.ServerAddress": &c.Worker.EventWriters.RPC.ServerAddress,
    "Worker.EventWriters.RPC.ServerPassword": &c.Worker.EventWriters.RPC.ServerPassword,
    "Worker.EventWriters.RPC.UpdateTimeout": &c.Worker.EventWriters.RPC.UpdateTimeout,
    "Worker.EventWriters.Dynamo.Region": &c.Worker.EventWriters.Dynamo.Region,
    "Worker.EventWriters.Dynamo.Key": &c.Worker.EventWriters.Dynamo.Key,
    "Worker.EventWriters.Dynamo.Secret": &c.Worker.EventWriters.Dynamo.Secret,
    "Worker.EventWriters.Dynamo.TableBasename": &c.Worker.EventWriters.Dynamo.TableBasename,
    "Scheduler.ScheduleRate": &c.Scheduler.ScheduleRate,
    "Scheduler.ScheduleChunk": &c.Scheduler.ScheduleChunk,
    "Scheduler.NodePingTimeout": &c.Scheduler.NodePingTimeout,
    "Scheduler.NodeInitTimeout": &c.Scheduler.NodeInitTimeout,
    "Scheduler.NodeDeadTimeout": &c.Scheduler.NodeDeadTimeout,
    "Scheduler.Node.ID": &c.Scheduler.Node.ID,
    "Scheduler.Node.Resources.Cpus": &c.Scheduler.Node.Resources.Cpus,
    "Scheduler.Node.Resources.RamGb": &c.Scheduler.Node.Resources.RamGb,
    "Scheduler.Node.Resources.DiskGb": &c.Scheduler.Node.Resources.DiskGb,
    "Scheduler.Node.Timeout": &c.Scheduler.Node.Timeout,
    "Scheduler.Node.UpdateRate": &c.Scheduler.Node.UpdateRate,
    "Scheduler.Node.UpdateTimeout": &c.Scheduler.Node.UpdateTimeout,
    "Scheduler.Node.Metadata": &c.Scheduler.Node.Metadata,
    "Scheduler.Node.ServerAddress": &c.Scheduler.Node.ServerAddress,
    "Scheduler.Node.ServerPassword": &c.Scheduler.Node.ServerPassword,
    "Scheduler.Node.Logger.Level": &c.Scheduler.Node.Logger.Level,
    "Scheduler.Node.Logger.Formatter": &c.Scheduler.Node.Logger.Formatter,
    "Scheduler.Node.Logger.OutputFile": &c.Scheduler.Node.Logger.OutputFile,
    "Scheduler.Node.Logger.Foo.FooField": &c.Scheduler.Node.Logger.Foo.FooField,
    "Scheduler.Node.Logger.Foo.Level": &c.Scheduler.Node.Logger.Foo.Level,
    "Scheduler.Logger.Level": &c.Scheduler.Logger.Level,
    "Scheduler.Logger.Formatter": &c.Scheduler.Logger.Formatter,
    "Scheduler.Logger.OutputFile": &c.Scheduler.Logger.OutputFile,
    "Scheduler.Logger.Foo.FooField": &c.Scheduler.Logger.Foo.FooField,
    "Scheduler.Logger.Foo.Level": &c.Scheduler.Logger.Foo.Level,
    "Scheduler.Worker.Storage.Local.AllowedDirs": &c.Scheduler.Worker.Storage.Local.AllowedDirs,
    "Scheduler.Worker.Storage.S3.Key": &c.Scheduler.Worker.Storage.S3.Key,
    "Scheduler.Worker.Storage.S3.Secret": &c.Scheduler.Worker.Storage.S3.Secret,
    "Scheduler.Worker.Storage.S3.FromEnv": &c.Scheduler.Worker.Storage.S3.FromEnv,
    "Scheduler.Worker.Storage.GS": &c.Scheduler.Worker.Storage.GS,
    "Scheduler.Worker.Storage.Swift.UserName": &c.Scheduler.Worker.Storage.Swift.UserName,
    "Scheduler.Worker.Storage.Swift.Password": &c.Scheduler.Worker.Storage.Swift.Password,
    "Scheduler.Worker.Storage.Swift.AuthURL": &c.Scheduler.Worker.Storage.Swift.AuthURL,
    "Scheduler.Worker.Storage.Swift.TenantName": &c.Scheduler.Worker.Storage.Swift.TenantName,
    "Scheduler.Worker.Storage.Swift.TenantID": &c.Scheduler.Worker.Storage.Swift.TenantID,
    "Scheduler.Worker.Storage.Swift.RegionName": &c.Scheduler.Worker.Storage.Swift.RegionName,
    "Scheduler.Worker.WorkDir": &c.Scheduler.Worker.WorkDir,
    "Scheduler.Worker.UpdateRate": &c.Scheduler.Worker.UpdateRate,
    "Scheduler.Worker.BufferSize": &c.Scheduler.Worker.BufferSize,
    "Scheduler.Worker.Logger.Level": &c.Scheduler.Worker.Logger.Level,
    "Scheduler.Worker.Logger.Formatter": &c.Scheduler.Worker.Logger.Formatter,
    "Scheduler.Worker.Logger.OutputFile": &c.Scheduler.Worker.Logger.OutputFile,
    "Scheduler.Worker.Logger.Foo.FooField": &c.Scheduler.Worker.Logger.Foo.FooField,
    "Scheduler.Worker.Logger.Foo.Level": &c.Scheduler.Worker.Logger.Foo.Level,
    "Scheduler.Worker.TaskReader": &c.Scheduler.Worker.TaskReader,
    "Scheduler.Worker.TaskReaders.RPC.ServerAddress": &c.Scheduler.Worker.TaskReaders.RPC.ServerAddress,
    "Scheduler.Worker.TaskReaders.RPC.ServerPassword": &c.Scheduler.Worker.TaskReaders.RPC.ServerPassword,
    "Scheduler.Worker.TaskReaders.Dynamo.Region": &c.Scheduler.Worker.TaskReaders.Dynamo.Region,
    "Scheduler.Worker.TaskReaders.Dynamo.Key": &c.Scheduler.Worker.TaskReaders.Dynamo.Key,
    "Scheduler.Worker.TaskReaders.Dynamo.Secret": &c.Scheduler.Worker.TaskReaders.Dynamo.Secret,
    "Scheduler.Worker.TaskReaders.Dynamo.TableBasename": &c.Scheduler.Worker.TaskReaders.Dynamo.TableBasename,
    "Scheduler.Worker.ActiveEventWriters": &c.Scheduler.Worker.ActiveEventWriters,
    "Scheduler.Worker.EventWriters.RPC.ServerAddress": &c.Scheduler.Worker.EventWriters.RPC.ServerAddress,
    "Scheduler.Worker.EventWriters.RPC.ServerPassword": &c.Scheduler.Worker.EventWriters.RPC.ServerPassword,
    "Scheduler.Worker.EventWriters.RPC.UpdateTimeout": &c.Scheduler.Worker.EventWriters.RPC.UpdateTimeout,
    "Scheduler.Worker.EventWriters.Dynamo.Region": &c.Scheduler.Worker.EventWriters.Dynamo.Region,
    "Scheduler.Worker.EventWriters.Dynamo.Key": &c.Scheduler.Worker.EventWriters.Dynamo.Key,
    "Scheduler.Worker.EventWriters.Dynamo.Secret": &c.Scheduler.Worker.EventWriters.Dynamo.Secret,
    "Scheduler.Worker.EventWriters.Dynamo.TableBasename": &c.Scheduler.Worker.EventWriters.Dynamo.TableBasename,
    "Scheduler.Backend": &c.Scheduler.Backend,
    "Scheduler.Backends.Local": &c.Scheduler.Backends.Local,
    "Scheduler.Backends.HTCondor.Template": &c.Scheduler.Backends.HTCondor.Template,
    "Scheduler.Backends.SLURM.Template": &c.Scheduler.Backends.SLURM.Template,
    "Scheduler.Backends.PBS.Template": &c.Scheduler.Backends.PBS.Template,
    "Scheduler.Backends.GridEngine.Template": &c.Scheduler.Backends.GridEngine.Template,
    "Scheduler.Backends.OpenStack.KeyPair": &c.Scheduler.Backends.OpenStack.KeyPair,
    "Scheduler.Backends.OpenStack.ConfigPath": &c.Scheduler.Backends.OpenStack.ConfigPath,
    "Scheduler.Backends.GCE.AccountFile": &c.Scheduler.Backends.GCE.AccountFile,
    "Scheduler.Backends.GCE.Project": &c.Scheduler.Backends.GCE.Project,
    "Scheduler.Backends.GCE.Zone": &c.Scheduler.Backends.GCE.Zone,
    "Scheduler.Backends.GCE.Weights.PreferQuickStartup": &c.Scheduler.Backends.GCE.Weights.PreferQuickStartup,
    "Scheduler.Backends.GCE.CacheTTL": &c.Scheduler.Backends.GCE.CacheTTL,
    "Log.Level": &c.Log.Level,
    "Log.Formatter": &c.Log.Formatter,
    "Log.OutputFile": &c.Log.OutputFile,
    "Log.Foo.FooField": &c.Log.Foo.FooField,
    "Log.Foo.Level": &c.Log.Foo.Level,
    "Dynamo.Region": &c.Dynamo.Region,
    "Dynamo.Key": &c.Dynamo.Key,
    "Dynamo.Secret": &c.Dynamo.Secret,
    "Dynamo.TableBasename": &c.Dynamo.TableBasename,
  }
}

func (c *Config) Set(k string, v interface{}) error {
}

func (c *Config) FlagVal(key string) flag.Value {
  return nil
}

func (c *Config) FlagSet() *flag.FlagSet {
  return nil
}

type flagval struct {
  ptr interface{}
}

func (f *flagval) String() string {
  return ""
}

func (f *flagval) CoerceSet(val interface{}) error {
  return roger.CoerceSet(f.ptr, v)
}

func (f *flagval) Set(s string) error {
  return nil
}

func (f *flagval) Get() interface{} {
  return nil
}
