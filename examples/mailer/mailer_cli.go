package main

import (
	"github.com/buchanae/cli"
	"github.com/buchanae/cli/examples/mailer/foo"
	"github.com/buchanae/mailer/imap"
	"github.com/buchanae/mailer/model"
	"github.com/sanity-io/litter"
)

func main() {
	p := cli.Providers(
		cli.Env("MAILER"),
		cli.YAMLFile("config.yaml"),
	)
	b := cli.Cobra{}
	b.Command.Use = "mailer"

	for _, spec := range cmdSpecs {
		b.AddSpec(spec, p)
	}

	b.Execute()
}

func initDB(opt Opt) *model.DB {
	db, err := model.Open(opt.DB.Path)
	cli.Check(err)
	return db
}

// Create a mailbox.
//
// Create a new mailbox in the database.
//
// Usage: mailer create mailbox <mailbox name>
// Example: mailer create mailbox foobar
func CreateMailbox(opt Opt, name string) {
	db := initDB(opt)
	defer db.Close()
	cli.Check(db.CreateMailbox(name))
}

func DeleteMailbox(opt Opt, name string) {
	db := initDB(opt)
	defer db.Close()
	cli.Check(db.DeleteMailbox(name))
}

func RenameMailbox(opt Opt, from, to string) {
	db := initDB(opt)
	defer db.Close()
	cli.Check(db.RenameMailbox(from, to))
}

// TODO have framework handle lots of init and coordiation?
//      or just keep it simple?

func GetMessage(opt Opt, ids ...int) {
	db := initDB(opt)
	defer db.Close()

	for _, id := range ids {
		msg, err := db.Message(id)
		cli.Check(err)
		litter.Dump(msg)
	}
}

func CreateMessage(opt Opt, mailbox, path string) {
	db := initDB(opt)
	defer db.Close()

	fh := cli.Open(path)
	defer fh.Close()

	_, err := db.CreateMessage(mailbox, fh, []imap.Flag{imap.Recent})
	cli.Check(err)
}

func ListMailboxes(opt Opt) {
	db := initDB(opt)
	defer db.Close()

	boxes, err := db.ListMailboxes()
	cli.Check(err)

	for _, box := range boxes {
		litter.Dump(box)
	}
}

func Foo(opt foo.Config) {
}

func Noarg() {
}
