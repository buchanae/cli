package cli

import (
  "strings"
  "github.com/spf13/cobra"
)

type Cobra struct {
  cobra.Command
  sub map[string]*cobra.Command
}

func (cb *Cobra) AddSpec(spec CmdSpec, p Provider) *cobra.Command {
  cmd := cb.addCobra(&cb.Command, 0, spec)
  fp := &PFlagProvider{FlagSet: cmd.Flags()}
  fp.AddOpts(spec.OptSpecs())
  p = Providers(p, fp)

  cmd.RunE = func(_ *cobra.Command, args []string) error {
    err := p.Init()
    if err != nil {
      return err
    }
    LoadOpts(p, spec.OptSpecs())
    return Run(spec, args)
  }
  return cmd
}

func (cb *Cobra) addCobra(cmd *cobra.Command, depth int, spec CmdSpec) *cobra.Command {
  if cb.sub == nil {
    cb.sub = map[string]*cobra.Command{}
  }

  det := ParseCmdDetail(spec)

  if depth == len(det.Path) - 1 {
    x := &cobra.Command{
      Use: det.Name,
      Short: det.Synopsis,
      Long: det.Doc,
      Example: det.Example,
      Deprecated: det.Deprecated,
      Hidden: det.Hidden,
      Aliases: det.Aliases,
    }
    cmd.AddCommand(x)
    return x
  }

  name := strings.Join(det.Path[:depth+1], " ")
  parent, ok := cb.sub[name]
  if !ok {
    parent = &cobra.Command{
      Use: name,
    }
    cmd.AddCommand(parent)
    cb.sub[name] = parent
  }
  return cb.addCobra(parent, depth + 1, spec)
}
