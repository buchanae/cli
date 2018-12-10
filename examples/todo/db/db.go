package db

import (
  "encoding/json"
  "fmt"
  "os"
  "time"
)

type DB struct {
  path string
}

func (db *DB) Add(description string, snooze time.Duration) (Todo, error) {
  data, err := load(db.path)
  if err != nil {
    return Todo{}, err
  }
  todo := Todo{
    ID: data.NextID,
    Description: description,
    Due: time.Now().Add(snooze),
  }
  data.NextID++
  data.Todos = append(data.Todos, todo)
  err = save(db.path, data)
  if err != nil {
    return Todo{}, err
  }
  return todo, nil
}

func (db *DB) Get(id int) (t Todo, err error) {
  data, err := load(db.path)
  if err != nil {
    return
  }

  for _, todo := range data.Todos {
    if todo.ID == id {
      t = todo
      return
    }
  }
  err = fmt.Errorf("no todo with ID %d", id)
  return
}

func (db *DB) Update(t Todo) error {
  data, err := load(db.path)
  if err != nil {
    return err
  }

  var updated []Todo
  for _, todo := range data.Todos {
    if todo.ID == t.ID {
      todo = t
    }
    updated = append(updated, todo)
  }

  data.Todos = updated
  return save(db.path, data)
}

func (db *DB) Delete(id int) error {
  data, err := load(db.path)
  if err != nil {
    return err
  }

  var updated []Todo
  for _, todo := range data.Todos {
    if todo.ID == id {
      continue
    }
    updated = append(updated, todo)
  }

  data.Todos = updated
  return save(db.path, data)
}

func (db *DB) List() ([]Todo, error) {
  data, err := load(db.path)
  if err != nil {
    return nil, err
  }
  return data.Todos, nil
}

func load(path string) (*todoList, error) {
  fh, err := os.Open(path)
  if err != nil {
    return nil, fmt.Errorf("loading database: %v", err)
  }

  todos := &todoList{}
  dec := json.NewDecoder(fh)
  err = dec.Decode(todos)
  if err != nil {
    return nil, fmt.Errorf("decoding database: %v", err)
  }
  return todos, nil
}

func save(path string, data *todoList) error {
  b, err := json.Marshal(data)
  if err != nil {
    return fmt.Errorf("encoding database: %v", err)
  }

  fh, err := os.Create(path)
  if err != nil {
    return fmt.Errorf("saving database: %v", err)
  }
  defer fh.Close()

  _, err = fh.Write(b)
  if err != nil {
    return fmt.Errorf("saving database: %v", err)
  }
  return nil
}

type Todo struct {
  ID int
  Description string
  Due time.Time
}

type todoList struct {
  NextID int
  Todos []Todo
}

type Opt struct {
  Path string
}

func Open(opt Opt) (*DB, error) {
  if opt.Path == "" {
    return nil, fmt.Errorf("db path is empty")
  }

  s, err := os.Stat(opt.Path)
  if err != nil {
    if os.IsNotExist(err) {
      err := save(opt.Path, &todoList{NextID: 1})
      if err != nil {
        return nil, fmt.Errorf("creating database file at %q: %v", opt.Path, err)
      }
    } else {
      return nil, fmt.Errorf("checking database file: %v", err)
    }
  } else if !s.Mode().IsRegular() {
    return nil, fmt.Errorf("db path %q already exists but is not a regular file", opt.Path)
  }

  return &DB{path: opt.Path}, nil
}
