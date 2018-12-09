package main

import (
  "github.com/buchanae/mailer/model"
  "github.com/buchanae/mailer/imap"
  "github.com/buchanae/cli/examples/mailer/foo"
  "github.com/buchanae/cli"
  "github.com/sanity-io/litter"
)

func main() {
}

type DBOpt struct {
  // Path to database directory
  Path string
}

type Opt struct {
  DB DBOpt
  Foo foo.Config
}

func DefaultOpt() Opt {
  return Opt{
    DB: DBOpt{
      Path: "mailer.db",
    },
  }
}

/* TODO allow this form
var DefaultOpt = Opt{
  DB: DBOpt{
    Path: "mailer.db",
  },
}
*/

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
func CreateMailboxCmd(opt Opt, name string) {
  db := initDB(opt)
  defer db.Close()
  cli.Check(db.CreateMailbox(name))
}

func DeleteMailboxCmd(opt Opt, name string) {
  db := initDB(opt)
  defer db.Close()
  cli.Check(db.DeleteMailbox(name))
}

func RenameMailboxCmd(opt Opt, from, to string) {
  db := initDB(opt)
  defer db.Close()
  cli.Check(db.RenameMailbox(from, to))
}

// TODO have framework handle lots of init and coordiation?
//      or just keep it simple?

func GetMessageCmd(opt Opt, ids ...int) {
  db := initDB(opt)
  defer db.Close()

  for _, id := range ids {
    msg, err := db.Message(id)
    cli.Check(err)
    litter.Dump(msg)
  }
}

func CreateMessageCmd(opt Opt, mailbox, path string) {
  db := initDB(opt)
  defer db.Close()

  fh := cli.Open(path)
  defer fh.Close()

  _, err := db.CreateMessage(mailbox, fh, []imap.Flag{imap.Recent})
  cli.Check(err)
}

func ListMailboxesCmd(opt Opt) {
  db := initDB(opt)
  defer db.Close()

  boxes, err := db.ListMailboxes()
  cli.Check(err)

  for _, box := range boxes {
    litter.Dump(box)
  }
}

func FooCmd(opt foo.Config) {
}

func NoArgCmd() {
}
