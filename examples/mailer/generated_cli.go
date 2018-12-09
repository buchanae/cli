package main

import cli "github.com/buchanae/cli"
import foo "github.com/buchanae/cli/examples/mailer/foo"

var cmdSpecs = []cli.CmdSpec{

	&createMailboxCmdSpec{
		Opt: DefaultOpt(),
	},

	&deleteMailboxCmdSpec{
		Opt: DefaultOpt(),
	},

	&renameMailboxCmdSpec{
		Opt: DefaultOpt(),
	},

	&getMessageCmdSpec{
		Opt: DefaultOpt(),
	},

	&createMessageCmdSpec{
		Opt: DefaultOpt(),
	},

	&listMailboxesCmdSpec{
		Opt: DefaultOpt(),
	},

	&fooCmdSpec{
		Opt: foo.DefaultConfig(),
	},

	&noArgCmdSpec{},
}

type createMailboxCmdSpec struct {
	Opt Opt
}

func (cmd *createMailboxCmdSpec) Name() string {
	return "CreateMailboxCmd"
}

func (cmd *createMailboxCmdSpec) Doc() string {
	return "Create a mailbox.\n\nCreate a new mailbox in the database.\n\nUsage: mailer create mailbox <mailbox name>\nExample: mailer create mailbox foobar\n"
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

type deleteMailboxCmdSpec struct {
	Opt Opt
}

func (cmd *deleteMailboxCmdSpec) Name() string {
	return "DeleteMailboxCmd"
}

func (cmd *deleteMailboxCmdSpec) Doc() string {
	return ""
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

type renameMailboxCmdSpec struct {
	Opt Opt
}

func (cmd *renameMailboxCmdSpec) Name() string {
	return "RenameMailboxCmd"
}

func (cmd *renameMailboxCmdSpec) Doc() string {
	return ""
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

type getMessageCmdSpec struct {
	Opt Opt
}

func (cmd *getMessageCmdSpec) Name() string {
	return "GetMessageCmd"
}

func (cmd *getMessageCmdSpec) Doc() string {
	return ""
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

type createMessageCmdSpec struct {
	Opt Opt
}

func (cmd *createMessageCmdSpec) Name() string {
	return "CreateMessageCmd"
}

func (cmd *createMessageCmdSpec) Doc() string {
	return ""
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

type listMailboxesCmdSpec struct {
	Opt Opt
}

func (cmd *listMailboxesCmdSpec) Name() string {
	return "ListMailboxesCmd"
}

func (cmd *listMailboxesCmdSpec) Doc() string {
	return ""
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

type fooCmdSpec struct {
	Opt foo.Config
}

func (cmd *fooCmdSpec) Name() string {
	return "FooCmd"
}

func (cmd *fooCmdSpec) Doc() string {
	return ""
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

type noArgCmdSpec struct {
}

func (cmd *noArgCmdSpec) Name() string {
	return "NoArgCmd"
}

func (cmd *noArgCmdSpec) Doc() string {
	return ""
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

