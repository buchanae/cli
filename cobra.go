package cli

import (
  "reflect"
  "fmt"
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
func (cb *Cobra) Add(spec Spec, l *Loader) *cobra.Command {

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

  f := x.Flags()
  for _, opt := range cmd.Opts {
    k := cb.keyfunc(opt.Key)
    f.VarP(&pflagOpt{opt}, k, opt.Short, opt.Synopsis)
    if opt.Deprecated != "" {
      f.MarkDeprecated(k, opt.Deprecated)
    }
    if opt.Hidden {
      f.MarkHidden(k)
    }
  }

  x.RunE = func(_ *cobra.Command, args []string) error {
    return Run(spec, l, args)
  }

  parent.AddCommand(x)
  return x
}

func (cb *Cobra) keyfunc(key []string) string {
  if cb.KeyFunc == nil {
    return DotKey(key)
  }
  return cb.KeyFunc(key)
}

type pflagOpt struct {
  opt *Opt
}

func (p *pflagOpt) Set(v string) error {
  return nil
}

func (p *pflagOpt) String() string {
  // TODO not super happy about including reflect just for this.
  //      also not happy with the default doc, want "os.Stdout" for io.Writer
  //      in some cases.
  t := reflect.TypeOf(p.opt.DefaultValue)
  if t.Kind() == reflect.Ptr {
    return ""
  }
  if p.opt.DefaultValue == nil {
    return ""
  }
  return fmt.Sprintf("%v", p.opt.DefaultValue)
}

func (p *pflagOpt) Type() string {
  return p.opt.Type
}
