package cli

import (
  "time"
  "strings"
  "github.com/spf13/cobra"
  "github.com/spf13/pflag"
)

func AddPFlags(fs *pflag.FlagSet, specs []OptSpec) {
  for _, opt := range specs {
    det := ParseOptDetail(opt)
    k := strings.ToLower(strings.Join(opt.Key, "-"))

    switch z := opt.Value.(type) {
		case *uint:
			fs.UintVar(z, k, *z, det.Synopsis)
		case *uint8:
			fs.Uint8Var(z, k, *z, det.Synopsis)
		case *uint16:
			fs.Uint16Var(z, k, *z, det.Synopsis)
		case *uint32:
			fs.Uint32Var(z, k, *z, det.Synopsis)
		case *uint64:
			fs.Uint64Var(z, k, *z, det.Synopsis)
		case *int:
			fs.IntVar(z, k, *z, det.Synopsis)
		case *int8:
			fs.Int8Var(z, k, *z, det.Synopsis)
		case *int16:
			fs.Int16Var(z, k, *z, det.Synopsis)
		case *int32:
			fs.Int32Var(z, k, *z, det.Synopsis)
		case *int64:
			fs.Int64Var(z, k, *z, det.Synopsis)
		case *float32:
			fs.Float32Var(z, k, *z, det.Synopsis)
		case *float64:
			fs.Float64Var(z, k, *z, det.Synopsis)
		case *bool:
			fs.BoolVar(z, k, *z, det.Synopsis)
		case *string:
			fs.StringVar(z, k, *z, det.Synopsis)
		case *[]string:
      fs.StringSliceVar(z, k, *z, det.Synopsis)
		case *time.Duration:
			fs.DurationVar(z, k, *z, det.Synopsis)
    default:
      // TODO should probably return error
    }
  }
}

func AddCobra(cmd *cobra.Command, specs ...CmdSpec) {
  addCobra(cmd, 0, specs...)
}

func addCobra(cmd *cobra.Command, depth int, specs ...CmdSpec) {
  sub := map[string]*cobra.Command{}

  for _, spec := range specs {
    det := ParseCmdDetail(spec)

    if depth == len(det.Path) - 1 {
      addCobraDet(cmd, det, spec)
      continue
    }

    name := det.Path[depth]
    parent, ok := sub[name]
    if !ok {
      parent = &cobra.Command{
        Use: name,
      }
      cmd.AddCommand(parent)
      sub[name] = parent
    }
    addCobra(parent, depth + 1, spec)
  }
}

func addCobraDet(cmd *cobra.Command, det CmdDetail, spec CmdSpec) {
  c := &cobra.Command{
    Use: det.Name,
    Short: det.Synopsis,
    Long: det.Doc,
    Example: det.Example,
    Deprecated: det.Deprecated,
    Hidden: det.Hidden,
    Aliases: det.Aliases,
    RunE: func(cmd *cobra.Command, args []string) error {
      return Run(spec, args)
    },
  }
  cmd.AddCommand(c)
  AddPFlags(c.Flags(), spec.OptSpecs())
}
