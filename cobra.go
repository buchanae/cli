package cli

import (
	"github.com/spf13/cobra"
)

// Cobra helps build a set of cobra commands.
// Commands are into a tree of subcommands based on their common subpaths.
type Cobra struct {
	cobra.Command
  KeyFunc
}

// Add adds a command to the tree.
// Option values are loaded from the given provider when the command runs.
// Commands are into a tree of subcommands based on their common subpaths.
func (cb *Cobra) Add(spec Spec) *cobra.Command {

  cmd := spec.Cmd()
  // TODO what's the right place for this?
  Enrich(cmd)
  parent, missing, _ := cb.Command.Find(cmd.Path)

  // Add missing intermediate commands.
  for i := 0; i < len(missing) - 1; i++ {
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

func (cb *Cobra) SetRunner(cmd *cobra.Command, spec Spec, l *Loader) {
  cmd.RunE = func(_ *cobra.Command, args []string) error {
    return Run(spec, l, args)
  }
}
