roger \
  -i Scheduler.Backends \
  -i Scheduler.Worker \
  -i Worker.EventWriters.Dynamo \
  -i Worker.Storage \
  -i Worker.EventWriters.RPC \
  -i Worker.TaskReaders.RPC \
  -i Worker.TaskReaders.Dynamo \
  -a w=Worker.WorkDir \
  -a host=Server.HostName \
  -root Config \
  -out ./example/gen.go \
  ./example/
