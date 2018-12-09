package cli

type CmdSpec interface {
  Name() string
  Doc() string
  Run(args []string)
  ArgSpecs() []ArgSpec
  OptSpecs() []OptSpec
}

type OptSpec struct {
  Key []string
  Doc string
  Type string
  Value interface{}
}

type ArgSpec struct {
  Name string
  Type string
  Variadic bool
}
