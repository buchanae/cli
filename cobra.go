package cli

import (
	"github.com/spf13/cobra"
)

// Cobra helps build a set of cobra commands.
// Commands are built into a tree of subcommands based on their common subpaths.
type Cobra struct {
	cobra.Command
	KeyFunc
}

// Add adds a command to the tree.
func (cb *Cobra) Add(spec Spec) *cobra.Command {

	cmd := spec.Cmd()
	parent, missing, _ := cb.Command.Find(cmd.Path)

	// Add missing intermediate commands.
	for i := 0; i < len(missing)-1; i++ {
		z := &cobra.Command{
			Use: missing[i],
		}
		parent.AddCommand(z)
		parent = z
	}

	x := &cobra.Command{
		Use:        cmd.Name,
		Short:      cmd.Synopsis,
		Long:       cmd.Doc,
		Example:    cmd.Example,
		Deprecated: cmd.Deprecated,
		Hidden:     cmd.Hidden,
		Aliases:    cmd.Aliases,
	}

	parent.AddCommand(x)
	return x
}

// SetRunner sets `cobra.Command.RunE` to use the loader and runner
// from this package.
func (cb *Cobra) SetRunner(cmd *cobra.Command, spec Spec, l *Loader) {
	cmd.RunE = func(_ *cobra.Command, args []string) error {
		return Run(spec, l, args)
	}
}
