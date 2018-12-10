package main

import (
  "fmt"
  "time"
  "github.com/buchanae/cli"
  "github.com/buchanae/cli/examples/todo/db"
)

func main() {
  p := cli.Providers(
    cli.Env("TODO"),
    cli.YAMLFile("config.yaml"),
  )
  b := cli.Cobra{}
  b.Command.Use = "todo"
  b.SilenceUsage = true

  for _, spec := range cmdSpecs {
    b.AddSpec(spec, p)
  }

  b.Execute()
}

// Add a new todo item.
// Usage
// Example: todo add --snooze 5d "get a life!"
func Add(opt AddOpt, description string) {
  db := openDB(opt.Opt)
  todo, err := db.Add(description, opt.Snooze)
  cli.Check(err)
  fmt.Printf("Added todo #%d\n", todo.ID)
}

// Delete todo items.
// Aliases: del
// Example: todo delete 1 2
func Delete(opt Opt, ids ...int) {
  db := openDB(opt)
  for _, id := range ids {
    err := db.Delete(id)
    cli.Check(err)
  }
}

// List all todo items.
func List(opt Opt) {
  db := openDB(opt)
  list, err := db.List()
  cli.Check(err)
  for _, todo := range list {
    fmt.Fprintln(opt.Stdout, todo.Description, todo.Due)
  }
}

// Snooze a todo item.
// Aliases: snz
// Example: todo snooze 1 3h
func Snooze(opt Opt, id int, dur time.Duration) {
  db := openDB(opt)
  todo, err := db.Get(id)
  cli.Check(err)
  todo.Due = todo.Due.Add(dur)
  err = db.Update(todo)
  cli.Check(err)
}

// Remove a todo item.
// Deprecated: remove has been renamed to "delete".
// Hidden
// Example: todo remove 1
func Remove(opt Opt, id int) {
  Delete(opt, id)
}

func openDB(opt Opt) *db.DB {
  db, err := db.Open(opt.DB)
  cli.Check(err)
  return db
}
