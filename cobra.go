package cli

import (
	"github.com/spf13/cobra"
)

// Cobra helps build a set of cobra commands.
// Commands are into a tree of subcommands based on their common subpaths.
type Cobra struct {
	cobra.Command
}

// AddSpec adds a command to the tree.
// Option values are loaded from the given provider when the command runs.
// Commands are into a tree of subcommands based on their common subpaths.
func (cb *Cobra) AddSpec(spec CmdSpec, p Provider) *cobra.Command {

	det := ParseCmdDetail(spec)
  parent, missing, _ := cb.Command.Find(det.Path)

  // Add missing intermediate commands.
  for i := 0; i < len(missing) - 1; i++ {
    z := &cobra.Command{
			Use: missing[i],
		}
		parent.AddCommand(z)
    parent = z
  }

  x := &cobra.Command{
    Use:        det.Name,
    Short:      det.Synopsis,
    Long:       det.Doc,
    Example:    det.Example,
    Deprecated: det.Deprecated,
    Hidden:     det.Hidden,
    Aliases:    det.Aliases,
  }

  // Set up flags.
	fp := &PFlagProvider{FlagSet: x.Flags()}
	fp.AddOpts(spec.OptSpecs())
	p = Providers(p, fp)

  // Set runner.
  x.RunE = func(_ *cobra.Command, args []string) error {
    err := LoadOpts(spec.OptSpecs(), p)
    if err != nil {
      return err
    }
    return Run(spec, args)
  }

  parent.AddCommand(x)
  return x
}
