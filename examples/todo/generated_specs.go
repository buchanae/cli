package main

import cli "github.com/buchanae/cli"
import time "time"

var cmdSpecs = []cli.CmdSpec{
	&addSpec{
		Opt: DefaultAddOpt(),
	},
	&deleteSpec{
		Opt: DefaultOpt(),
	},
	&listSpec{
		Opt: DefaultOpt(),
	},
	&snoozeSpec{
		Opt: DefaultOpt(),
	},
	&removeSpec{
		Opt: DefaultOpt(),
	},
}

type addSpec struct {
	Opt AddOpt

	args struct {
		arg0 string
	}
}

func (cmd *addSpec) Name() string {
	return "Add"
}

func (cmd *addSpec) Doc() string {
	return "Add a new todo item.\nUsage\nExample: todo add --snooze 5d \"get a life!\"\n"
}

func (cmd *addSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	Add(
		cmd.Opt,
		cmd.args.arg0,
	)
}

func (cmd *addSpec) ArgSpecs() []cli.ArgSpec {
	return []cli.ArgSpec{
		{
			Name:     "description",
			Type:     "string",
			Variadic: false,
			Value:    &cmd.args.arg0,
		},
	}
}

func (cmd *addSpec) OptSpecs() []cli.OptSpec {
	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Stdout"},
			Doc:   "",
			Value: &cmd.Opt.Stdout,
		}, {
			Key:   []string{"Snooze"},
			Doc:   "",
			Value: &cmd.Opt.Snooze,
		},
	}
}

type deleteSpec struct {
	Opt Opt

	args struct {
		arg0 []int
	}
}

func (cmd *deleteSpec) Name() string {
	return "Delete"
}

func (cmd *deleteSpec) Doc() string {
	return "Delete todo items.\nAliases: del\nExample: todo delete 1 2\n"
}

func (cmd *deleteSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	Delete(
		cmd.Opt,
		cmd.args.arg0...,
	)
}

func (cmd *deleteSpec) ArgSpecs() []cli.ArgSpec {
	return []cli.ArgSpec{
		{
			Name:     "ids",
			Type:     "[]int",
			Variadic: true,
			Value:    &cmd.args.arg0,
		},
	}
}

func (cmd *deleteSpec) OptSpecs() []cli.OptSpec {
	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Stdout"},
			Doc:   "",
			Value: &cmd.Opt.Stdout,
		},
	}
}

type listSpec struct {
	Opt Opt

	args struct {
	}
}

func (cmd *listSpec) Name() string {
	return "List"
}

func (cmd *listSpec) Doc() string {
	return "List all todo items.\n"
}

func (cmd *listSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	List(
		cmd.Opt,
	)
}

func (cmd *listSpec) ArgSpecs() []cli.ArgSpec {

	return nil

}

func (cmd *listSpec) OptSpecs() []cli.OptSpec {
	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Stdout"},
			Doc:   "",
			Value: &cmd.Opt.Stdout,
		},
	}
}

type snoozeSpec struct {
	Opt Opt

	args struct {
		arg0 int

		arg1 time.Duration
	}
}

func (cmd *snoozeSpec) Name() string {
	return "Snooze"
}

func (cmd *snoozeSpec) Doc() string {
	return "Snooze a todo item.\nAliases: snz\nExample: todo snooze 1 3h\n"
}

func (cmd *snoozeSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	Snooze(
		cmd.Opt,
		cmd.args.arg0,
		cmd.args.arg1,
	)
}

func (cmd *snoozeSpec) ArgSpecs() []cli.ArgSpec {
	return []cli.ArgSpec{
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
	}
}

func (cmd *snoozeSpec) OptSpecs() []cli.OptSpec {
	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Stdout"},
			Doc:   "",
			Value: &cmd.Opt.Stdout,
		},
	}
}

type removeSpec struct {
	Opt Opt

	args struct {
		arg0 int
	}
}

func (cmd *removeSpec) Name() string {
	return "Remove"
}

func (cmd *removeSpec) Doc() string {
	return "Remove a todo item.\nDeprecated: remove has been renamed to \"delete\".\nHidden\nExample: todo remove 1\n"
}

func (cmd *removeSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	Remove(
		cmd.Opt,
		cmd.args.arg0,
	)
}

func (cmd *removeSpec) ArgSpecs() []cli.ArgSpec {
	return []cli.ArgSpec{
		{
			Name:     "id",
			Type:     "int",
			Variadic: false,
			Value:    &cmd.args.arg0,
		},
	}
}

func (cmd *removeSpec) OptSpecs() []cli.OptSpec {
	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Stdout"},
			Doc:   "",
			Value: &cmd.Opt.Stdout,
		},
	}
}

