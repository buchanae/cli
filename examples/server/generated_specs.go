package main

import cli "github.com/buchanae/cli"

func specs() []cli.Spec {
	return []cli.Spec{
		&defaultServerOptSpec{},
		&runSpec{
			opt: DefaultServerOpt(),
		},
	}
}

type defaultServerOptSpec struct {
	cmd *cli.Cmd

	args struct {
	}
}

func (cmd *defaultServerOptSpec) Run() {
	DefaultServerOpt()
}

func (cmd *defaultServerOptSpec) Cmd() *cli.Cmd {
	if cmd.cmd != nil {
		return cmd.cmd
	}
	cmd.cmd = &cli.Cmd{
		FuncName: "DefaultServerOpt",
		RawDoc:   "",
		Args:     []*cli.Arg{},
		Opts:     []*cli.Opt{},
	}
	return cmd.cmd
}

type runSpec struct {
	cmd  *cli.Cmd
	opt  ServerOpt
	args struct {
		arg0 string
	}
}

func (cmd *runSpec) Run() {
	Run(
		cmd.opt,
		cmd.args.arg0,
	)
}

func (cmd *runSpec) Cmd() *cli.Cmd {
	if cmd.cmd != nil {
		return cmd.cmd
	}
	cmd.cmd = &cli.Cmd{
		FuncName: "Run",
		RawDoc:   "",
		Args: []*cli.Arg{
			{
				Name:     "msg",
				Type:     "string",
				Variadic: false,
				Value:    &cmd.args.arg0,
			},
		},
		Opts: []*cli.Opt{
			{
				Key:          []string{"Name"},
				RawDoc:       "Server name, for metadata endpoints.\n",
				Value:        &cmd.opt.Name,
				DefaultValue: cmd.opt.Name,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"Addr"},
				RawDoc:       "Address to listen on.\n",
				Value:        &cmd.opt.Addr,
				DefaultValue: cmd.opt.Addr,
				Type:         "string",
				Short:        "",
			},
		},
	}
	return cmd.cmd
}
