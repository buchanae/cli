package main

import cli "github.com/buchanae/cli"
import foo "github.com/buchanae/cli/examples/mailer/foo"

func newCreateMailboxCmdSpec() *createMailboxCmdSpec {
	return &createMailboxCmdSpec{
		CmdSpec: cli.CmdSpec{
			Name: "CreateMailboxCmd",
			Doc:  "Create a mailbox.\n\nCreate a new mailbox in the database.\n\nUsage: mailer create mailbox <mailbox name>\nExample: mailer create mailbox foobar\n",
		},
		Opt: DefaultOpt(),
	}
}

type createMailboxCmdSpec struct {
	cli.CmdSpec
	Opt Opt
}

func (cmd *createMailboxCmdSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	CreateMailboxCmd(
		cmd.Opt,
		cli.CoerceString(args[0]),
	)
}

func (cmd *createMailboxCmdSpec) ArgSpecs() []cli.ArgSpec {

	return []cli.ArgSpec{
		{
			Name:     "name",
			Type:     "string",
			Variadic: false,
		},
	}

}

func (cmd *createMailboxCmdSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Type:  "string",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Type:  "int",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

func newDeleteMailboxCmdSpec() *deleteMailboxCmdSpec {
	return &deleteMailboxCmdSpec{
		CmdSpec: cli.CmdSpec{
			Name: "DeleteMailboxCmd",
			Doc:  "",
		},
		Opt: DefaultOpt(),
	}
}

type deleteMailboxCmdSpec struct {
	cli.CmdSpec
	Opt Opt
}

func (cmd *deleteMailboxCmdSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	DeleteMailboxCmd(
		cmd.Opt,
		cli.CoerceString(args[0]),
	)
}

func (cmd *deleteMailboxCmdSpec) ArgSpecs() []cli.ArgSpec {

	return []cli.ArgSpec{
		{
			Name:     "name",
			Type:     "string",
			Variadic: false,
		},
	}

}

func (cmd *deleteMailboxCmdSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Type:  "string",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Type:  "int",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

func newRenameMailboxCmdSpec() *renameMailboxCmdSpec {
	return &renameMailboxCmdSpec{
		CmdSpec: cli.CmdSpec{
			Name: "RenameMailboxCmd",
			Doc:  "",
		},
		Opt: DefaultOpt(),
	}
}

type renameMailboxCmdSpec struct {
	cli.CmdSpec
	Opt Opt
}

func (cmd *renameMailboxCmdSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	RenameMailboxCmd(
		cmd.Opt,
		cli.CoerceString(args[0]),
		cli.CoerceString(args[1]),
	)
}

func (cmd *renameMailboxCmdSpec) ArgSpecs() []cli.ArgSpec {

	return []cli.ArgSpec{
		{
			Name:     "from",
			Type:     "string",
			Variadic: false,
		}, {
			Name:     "to",
			Type:     "string",
			Variadic: false,
		},
	}

}

func (cmd *renameMailboxCmdSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Type:  "string",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Type:  "int",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

func newGetMessageCmdSpec() *getMessageCmdSpec {
	return &getMessageCmdSpec{
		CmdSpec: cli.CmdSpec{
			Name: "GetMessageCmd",
			Doc:  "",
		},
		Opt: DefaultOpt(),
	}
}

type getMessageCmdSpec struct {
	cli.CmdSpec
	Opt Opt
}

func (cmd *getMessageCmdSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	GetMessageCmd(
		cmd.Opt,
		cli.CoerceInts(args[0:])...,
	)
}

func (cmd *getMessageCmdSpec) ArgSpecs() []cli.ArgSpec {

	return []cli.ArgSpec{
		{
			Name:     "ids",
			Type:     "[]int",
			Variadic: true,
		},
	}

}

func (cmd *getMessageCmdSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Type:  "string",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Type:  "int",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

func newCreateMessageCmdSpec() *createMessageCmdSpec {
	return &createMessageCmdSpec{
		CmdSpec: cli.CmdSpec{
			Name: "CreateMessageCmd",
			Doc:  "",
		},
		Opt: DefaultOpt(),
	}
}

type createMessageCmdSpec struct {
	cli.CmdSpec
	Opt Opt
}

func (cmd *createMessageCmdSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	CreateMessageCmd(
		cmd.Opt,
		cli.CoerceString(args[0]),
		cli.CoerceString(args[1]),
	)
}

func (cmd *createMessageCmdSpec) ArgSpecs() []cli.ArgSpec {

	return []cli.ArgSpec{
		{
			Name:     "mailbox",
			Type:     "string",
			Variadic: false,
		}, {
			Name:     "path",
			Type:     "string",
			Variadic: false,
		},
	}

}

func (cmd *createMessageCmdSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Type:  "string",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Type:  "int",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

func newListMailboxesCmdSpec() *listMailboxesCmdSpec {
	return &listMailboxesCmdSpec{
		CmdSpec: cli.CmdSpec{
			Name: "ListMailboxesCmd",
			Doc:  "",
		},
		Opt: DefaultOpt(),
	}
}

type listMailboxesCmdSpec struct {
	cli.CmdSpec
	Opt Opt
}

func (cmd *listMailboxesCmdSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	ListMailboxesCmd(
		cmd.Opt,
	)
}

func (cmd *listMailboxesCmdSpec) ArgSpecs() []cli.ArgSpec {

	return nil

}

func (cmd *listMailboxesCmdSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Type:  "string",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Type:  "int",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Type:  "string",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

func newFooCmdSpec() *fooCmdSpec {
	return &fooCmdSpec{
		CmdSpec: cli.CmdSpec{
			Name: "FooCmd",
			Doc:  "",
		},
		Opt: foo.DefaultConfig(),
	}
}

type fooCmdSpec struct {
	cli.CmdSpec
	Opt foo.Config
}

func (cmd *fooCmdSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	FooCmd(
		cmd.Opt,
	)
}

func (cmd *fooCmdSpec) ArgSpecs() []cli.ArgSpec {

	return nil

}

func (cmd *fooCmdSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"Port"},
			Doc:   "Server port to listen on.\n",
			Type:  "int",
			Value: &cmd.Opt.Port,
		}, {
			Key:   []string{"Host"},
			Doc:   "Server host to listen on.\n",
			Type:  "string",
			Value: &cmd.Opt.Host,
		}, {
			Key:   []string{"User", "Username"},
			Doc:   "User name for login.\n",
			Type:  "string",
			Value: &cmd.Opt.User.Username,
		}, {
			Key:   []string{"User", "Password"},
			Doc:   "Password for login.\n",
			Type:  "string",
			Value: &cmd.Opt.User.Password,
		},
	}

}

func newNoArgCmdSpec() *noArgCmdSpec {
	return &noArgCmdSpec{
		CmdSpec: cli.CmdSpec{
			Name: "NoArgCmd",
			Doc:  "",
		},
	}
}

type noArgCmdSpec struct {
	cli.CmdSpec
}

func (cmd *noArgCmdSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	NoArgCmd()
}

func (cmd *noArgCmdSpec) ArgSpecs() []cli.ArgSpec {

	return nil

}

func (cmd *noArgCmdSpec) OptSpecs() []cli.OptSpec {

	return nil

}

