package main

import (
  "github.com/buchanae/cli"
  "github.com/buchanae/funnel/cmd/worker"
  "github.com/buchanae/funnel/cmd/task"
)

type Runner interface {
  Run()
}

type HasDocs interface {
  Docs()
}

func main() {
  c := cli.New()

  c.YAML("funnel.config.yml", "funnel.config.yaml")
  c.YAMLFlag("config", "c")
  c.ConsulFlag("consul-config")
  c.GCEMetaFlag("gce-config")
  c.OpenstackFlag("openstack-config")
  c.EtcdFlag("etcd-config")
  c.Env("funnel")
  c.Flags()

  c.Cmd("task create", task.DefaultCreate)
  c.Cmd("task list", task.DefaultList)
  c.Cmd("task get", task.DefaultGet)
  c.Cmd("task cancel", task.DefaultCancel)
  c.Cmd("worker run", worker.DefaultRun)

  err := c.Run()
  if err != nil {
  }
}

type CLI struct {}

func New() *CLI {
  return &CLI{}
}

func Cmd(r Runner) {
}

func YAML(path string, alternatives ...string) {
}

func Consul(addr string, alternatives ...string) {
}

func GCEMeta(addr string, alternatives ...string) {
}

func Openstack(addr string, alternatives ...string) {
}

func Etcd(addr string, alternatives ...string) {
}

func YAMLFlag(path string, alternatives ...string) {
}

func ConsulFlag(flag string, alternatives ...string) {
}

func GCEMetaFlag(flag string, alternatives ...string) {
}

func OpenstackFlag(flag string, alternatives ...string) {
}

func EtcdFlag(flag string, alternatives ...string) {
}
