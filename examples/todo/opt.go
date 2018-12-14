package main

import (
	"github.com/buchanae/cli/examples/todo/db"
	"io"
	"os"
	"time"
)

type AddOpt struct {
	Opt
	Snooze time.Duration `short:"s"`
	Tags   map[string]string
}

func DefaultAddOpt() AddOpt {
	return AddOpt{
		Opt:    DefaultOpt(),
		Snooze: time.Hour * 24,
	}
}

type Opt struct {
	// Path to config file.
	Config string
	DB     db.Opt
	Stdout io.Writer
}

func DefaultOpt() Opt {
	return Opt{
		Config: "config.yaml",
		DB: db.Opt{
			Path: "todo.db.json",
		},
		Stdout: os.Stdout,
	}
}
