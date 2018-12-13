package main

import cli "github.com/buchanae/cli"
import time "time"

func specs() []cli.Spec {
	return []cli.Spec{
		&addSpec{
			opt: DefaultAddOpt(),
		},
		&deleteSpec{
			opt: DefaultOpt(),
		},
		&listSpec{
			opt: DefaultOpt(),
		},
		&snoozeSpec{
			opt: DefaultOpt(),
		},
		&removeSpec{
			opt: DefaultOpt(),
		},
	}
}

type addSpec struct {
	opt  AddOpt
	args struct {
		arg0 string
	}
}

func (cmd *addSpec) Run() {
	Add(
		cmd.opt,
		cmd.args.arg0,
	)
}

func (cmd *addSpec) Cmd() *cli.Cmd {
	return &cli.Cmd{
		FuncName: "Add",
		RawDoc:   "Add a new todo item.\nUsage\nExample: todo add --snooze 5d \"get a life!\"\n",
		Args: []*cli.Arg{
			{
				Name:     "description",
				Type:     "string",
				Variadic: false,
				Value:    &cmd.args.arg0,
			},
		},
		Opts: []*cli.Opt{
			{
				Key:          []string{"Config"},
				RawDoc:       "Path to config file.\n",
				Value:        &cmd.opt.Config,
				DefaultValue: cmd.opt.Config,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"DB", "Path"},
				RawDoc:       "",
				Value:        &cmd.opt.DB.Path,
				DefaultValue: cmd.opt.DB.Path,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"Stdout"},
				RawDoc:       "",
				Value:        &cmd.opt.Stdout,
				DefaultValue: cmd.opt.Stdout,
				Type:         "io.Writer",
				Short:        "",
			}, {
				Key:          []string{"Snooze"},
				RawDoc:       "",
				Value:        &cmd.opt.Snooze,
				DefaultValue: cmd.opt.Snooze,
				Type:         "time.Duration",
				Short:        "s",
			}, {
				Key:          []string{"Tags"},
				RawDoc:       "",
				Value:        &cmd.opt.Tags,
				DefaultValue: cmd.opt.Tags,
				Type:         "map[string]string",
				Short:        "",
			},
		},
	}
}

type deleteSpec struct {
	opt  Opt
	args struct {
		arg0 []int
	}
}

func (cmd *deleteSpec) Run() {
	Delete(
		cmd.opt,
		cmd.args.arg0...,
	)
}

func (cmd *deleteSpec) Cmd() *cli.Cmd {
	return &cli.Cmd{
		FuncName: "Delete",
		RawDoc:   "Delete todo items.\nAliases: del\nExample: todo delete 1 2\n",
		Args: []*cli.Arg{
			{
				Name:     "ids",
				Type:     "[]int",
				Variadic: true,
				Value:    &cmd.args.arg0,
			},
		},
		Opts: []*cli.Opt{
			{
				Key:          []string{"Config"},
				RawDoc:       "Path to config file.\n",
				Value:        &cmd.opt.Config,
				DefaultValue: cmd.opt.Config,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"DB", "Path"},
				RawDoc:       "",
				Value:        &cmd.opt.DB.Path,
				DefaultValue: cmd.opt.DB.Path,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"Stdout"},
				RawDoc:       "",
				Value:        &cmd.opt.Stdout,
				DefaultValue: cmd.opt.Stdout,
				Type:         "io.Writer",
				Short:        "",
			},
		},
	}
}

type listSpec struct {
	opt  Opt
	args struct {
	}
}

func (cmd *listSpec) Run() {
	List(
		cmd.opt,
	)
}

func (cmd *listSpec) Cmd() *cli.Cmd {
	return &cli.Cmd{
		FuncName: "List",
		RawDoc:   "List all todo items.\n",
		Args:     []*cli.Arg{},
		Opts: []*cli.Opt{
			{
				Key:          []string{"Config"},
				RawDoc:       "Path to config file.\n",
				Value:        &cmd.opt.Config,
				DefaultValue: cmd.opt.Config,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"DB", "Path"},
				RawDoc:       "",
				Value:        &cmd.opt.DB.Path,
				DefaultValue: cmd.opt.DB.Path,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"Stdout"},
				RawDoc:       "",
				Value:        &cmd.opt.Stdout,
				DefaultValue: cmd.opt.Stdout,
				Type:         "io.Writer",
				Short:        "",
			},
		},
	}
}

type snoozeSpec struct {
	opt  Opt
	args struct {
		arg0 int
		arg1 time.Duration
	}
}

func (cmd *snoozeSpec) Run() {
	Snooze(
		cmd.opt,
		cmd.args.arg0,
		cmd.args.arg1,
	)
}

func (cmd *snoozeSpec) Cmd() *cli.Cmd {
	return &cli.Cmd{
		FuncName: "Snooze",
		RawDoc:   "Snooze a todo item.\nAliases: snz\nExample: todo snooze 1 3h\n",
		Args: []*cli.Arg{
			{
				Name:     "id",
				Type:     "int",
				Variadic: false,
				Value:    &cmd.args.arg0,
			}, {
				Name:     "dur",
				Type:     "time.Duration",
				Variadic: false,
				Value:    &cmd.args.arg1,
			},
		},
		Opts: []*cli.Opt{
			{
				Key:          []string{"Config"},
				RawDoc:       "Path to config file.\n",
				Value:        &cmd.opt.Config,
				DefaultValue: cmd.opt.Config,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"DB", "Path"},
				RawDoc:       "",
				Value:        &cmd.opt.DB.Path,
				DefaultValue: cmd.opt.DB.Path,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"Stdout"},
				RawDoc:       "",
				Value:        &cmd.opt.Stdout,
				DefaultValue: cmd.opt.Stdout,
				Type:         "io.Writer",
				Short:        "",
			},
		},
	}
}

type removeSpec struct {
	opt  Opt
	args struct {
		arg0 int
	}
}

func (cmd *removeSpec) Run() {
	Remove(
		cmd.opt,
		cmd.args.arg0,
	)
}

func (cmd *removeSpec) Cmd() *cli.Cmd {
	return &cli.Cmd{
		FuncName: "Remove",
		RawDoc:   "Remove a todo item.\nDeprecated: remove has been renamed to \"delete\".\nHidden\nExample: todo remove 1\n",
		Args: []*cli.Arg{
			{
				Name:     "id",
				Type:     "int",
				Variadic: false,
				Value:    &cmd.args.arg0,
			},
		},
		Opts: []*cli.Opt{
			{
				Key:          []string{"Config"},
				RawDoc:       "Path to config file.\n",
				Value:        &cmd.opt.Config,
				DefaultValue: cmd.opt.Config,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"DB", "Path"},
				RawDoc:       "",
				Value:        &cmd.opt.DB.Path,
				DefaultValue: cmd.opt.DB.Path,
				Type:         "string",
				Short:        "",
			}, {
				Key:          []string{"Stdout"},
				RawDoc:       "",
				Value:        &cmd.opt.Stdout,
				DefaultValue: cmd.opt.Stdout,
				Type:         "io.Writer",
				Short:        "",
			},
		},
	}
}

