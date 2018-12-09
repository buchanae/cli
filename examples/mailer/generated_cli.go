package main

import cli "github.com/buchanae/cli"
import foo "github.com/buchanae/cli/examples/mailer/foo"

var cmdSpecs = []cli.CmdSpec{

	&createMailboxSpec{
		Opt: DefaultOpt(),
	},

	&deleteMailboxSpec{
		Opt: DefaultOpt(),
	},

	&renameMailboxSpec{
		Opt: DefaultOpt(),
	},

	&getMessageSpec{
		Opt: DefaultOpt(),
	},

	&createMessageSpec{
		Opt: DefaultOpt(),
	},

	&listMailboxesSpec{
		Opt: DefaultOpt(),
	},

	&fooSpec{
		Opt: foo.DefaultConfig(),
	},

	&noargSpec{},
}

type createMailboxSpec struct {
	Opt Opt
}

func (cmd *createMailboxSpec) Name() string {
	return "CreateMailbox"
}

func (cmd *createMailboxSpec) Doc() string {
	return "Create a mailbox.\n\nCreate a new mailbox in the database.\n\nUsage: mailer create mailbox <mailbox name>\nExample: mailer create mailbox foobar\n"
}

func (cmd *createMailboxSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	CreateMailboxCmd(
		cmd.Opt,
		cli.CoerceString(args[0]),
	)
}

func (cmd *createMailboxSpec) ArgSpecs() []cli.ArgSpec {

	return []cli.ArgSpec{
		{
			Name:     "name",
			Type:     "string",
			Variadic: false,
		},
	}

}

func (cmd *createMailboxSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

type deleteMailboxSpec struct {
	Opt Opt
}

func (cmd *deleteMailboxSpec) Name() string {
	return "DeleteMailbox"
}

func (cmd *deleteMailboxSpec) Doc() string {
	return ""
}

func (cmd *deleteMailboxSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	DeleteMailboxCmd(
		cmd.Opt,
		cli.CoerceString(args[0]),
	)
}

func (cmd *deleteMailboxSpec) ArgSpecs() []cli.ArgSpec {

	return []cli.ArgSpec{
		{
			Name:     "name",
			Type:     "string",
			Variadic: false,
		},
	}

}

func (cmd *deleteMailboxSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

type renameMailboxSpec struct {
	Opt Opt
}

func (cmd *renameMailboxSpec) Name() string {
	return "RenameMailbox"
}

func (cmd *renameMailboxSpec) Doc() string {
	return ""
}

func (cmd *renameMailboxSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	RenameMailboxCmd(
		cmd.Opt,
		cli.CoerceString(args[0]),
		cli.CoerceString(args[1]),
	)
}

func (cmd *renameMailboxSpec) ArgSpecs() []cli.ArgSpec {

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

func (cmd *renameMailboxSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

type getMessageSpec struct {
	Opt Opt
}

func (cmd *getMessageSpec) Name() string {
	return "GetMessage"
}

func (cmd *getMessageSpec) Doc() string {
	return ""
}

func (cmd *getMessageSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	GetMessageCmd(
		cmd.Opt,
		cli.CoerceInts(args[0:])...,
	)
}

func (cmd *getMessageSpec) ArgSpecs() []cli.ArgSpec {

	return []cli.ArgSpec{
		{
			Name:     "ids",
			Type:     "[]int",
			Variadic: true,
		},
	}

}

func (cmd *getMessageSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

type createMessageSpec struct {
	Opt Opt
}

func (cmd *createMessageSpec) Name() string {
	return "CreateMessage"
}

func (cmd *createMessageSpec) Doc() string {
	return ""
}

func (cmd *createMessageSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	CreateMessageCmd(
		cmd.Opt,
		cli.CoerceString(args[0]),
		cli.CoerceString(args[1]),
	)
}

func (cmd *createMessageSpec) ArgSpecs() []cli.ArgSpec {

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

func (cmd *createMessageSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

type listMailboxesSpec struct {
	Opt Opt
}

func (cmd *listMailboxesSpec) Name() string {
	return "ListMailboxes"
}

func (cmd *listMailboxesSpec) Doc() string {
	return ""
}

func (cmd *listMailboxesSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	ListMailboxesCmd(
		cmd.Opt,
	)
}

func (cmd *listMailboxesSpec) ArgSpecs() []cli.ArgSpec {

	return nil

}

func (cmd *listMailboxesSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"DB", "Path"},
			Doc:   "Path to database directory\n",
			Value: &cmd.Opt.DB.Path,
		}, {
			Key:   []string{"Foo", "Port"},
			Doc:   "Server port to listen on.\n",
			Value: &cmd.Opt.Foo.Port,
		}, {
			Key:   []string{"Foo", "Host"},
			Doc:   "Server host to listen on.\n",
			Value: &cmd.Opt.Foo.Host,
		}, {
			Key:   []string{"Foo", "User", "Username"},
			Doc:   "User name for login.\n",
			Value: &cmd.Opt.Foo.User.Username,
		}, {
			Key:   []string{"Foo", "User", "Password"},
			Doc:   "Password for login.\n",
			Value: &cmd.Opt.Foo.User.Password,
		},
	}

}

type fooSpec struct {
	Opt foo.Config
}

func (cmd *fooSpec) Name() string {
	return "Foo"
}

func (cmd *fooSpec) Doc() string {
	return ""
}

func (cmd *fooSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	FooCmd(
		cmd.Opt,
	)
}

func (cmd *fooSpec) ArgSpecs() []cli.ArgSpec {

	return nil

}

func (cmd *fooSpec) OptSpecs() []cli.OptSpec {

	return []cli.OptSpec{
		{
			Key:   []string{"Port"},
			Doc:   "Server port to listen on.\n",
			Value: &cmd.Opt.Port,
		}, {
			Key:   []string{"Host"},
			Doc:   "Server host to listen on.\n",
			Value: &cmd.Opt.Host,
		}, {
			Key:   []string{"User", "Username"},
			Doc:   "User name for login.\n",
			Value: &cmd.Opt.User.Username,
		}, {
			Key:   []string{"User", "Password"},
			Doc:   "Password for login.\n",
			Value: &cmd.Opt.User.Password,
		},
	}

}

type noargSpec struct {
}

func (cmd *noargSpec) Name() string {
	return "Noarg"
}

func (cmd *noargSpec) Doc() string {
	return ""
}

func (cmd *noargSpec) Run(args []string) {
	cli.CheckArgs(args, cmd.ArgSpecs())
	NoargCmd()
}

func (cmd *noargSpec) ArgSpecs() []cli.ArgSpec {

	return nil

}

func (cmd *noargSpec) OptSpecs() []cli.OptSpec {

	return nil

}

