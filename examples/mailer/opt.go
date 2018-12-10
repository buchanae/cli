package main

import (
	"github.com/buchanae/cli/examples/mailer/foo"
)

type DBOpt struct {
	// Path to database directory
	Path string
}

type Opt struct {
	DB  DBOpt
	Foo foo.Config
}

func DefaultOpt() Opt {
	return Opt{
		DB: DBOpt{
			Path: "mailer.data",
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
