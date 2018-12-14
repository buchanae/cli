package cli

/*
AutoCobra sets up a common pattern for apps:
commands are built into a tree of cobra subcommands,
flags are created using pflag, options are loaded
from env. vars, flags, and YAML config files.
*/
func AutoCobra(appname string, specs []Spec) error {
	b := Cobra{}
	b.Use = appname
	b.SilenceUsage = true

	for _, spec := range specs {
		cmd := b.Add(spec)
		opts := spec.Cmd().Opts
		flags := PFlags(cmd.Flags(), opts, DotKey)

		l := NewLoader(opts,
			Env(appname),
			flags,
			YAML(DefaultYAML),
		)
		b.SetRunner(cmd, spec, l)
	}

	return b.Execute()
}
