package cli

type CmdSpec struct {
  Name string
  Doc string
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
