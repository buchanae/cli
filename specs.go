package cli

// CmdSpec is implemented by types which define a CLI command.
// Normally these types are written by the CLI code generator.
type CmdSpec interface {
	Name() string
	Doc() string
	Run(args []string)
	ArgSpecs() []ArgSpec
	OptSpecs() []OptSpec
}

// OptSpec defines a config option for a CLI command.
// Normally instances are written by the CLI code generator.
type OptSpec struct {
	Key   []string
	Doc   string
	Value interface{}
}

// ArgSpec defines a positional argument for a CLI command.
// Normally instances are written by the CLI code generator.
type ArgSpec struct {
	Name     string
	Type     string
	Variadic bool
  Value interface{}
}
