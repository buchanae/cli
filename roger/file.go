package roger

import (
  "fmt"
  "path/filepath"
  "os"
)

type FileProvider struct {
  Keyfunc
  path string
  data map[string]interface{}
}

func NewFileProvider(path string) *FileProvider {
  return &FileProvider{path: path}
}

func OptionalFileProvider(path string) *FileProvider {
  if _, err := os.Stat(path); os.IsNotExist(err) {
    path = ""
  }
  return NewFileProvider(path)
}

func (f *FileProvider) Init() (err error) {
  if f.path != "" {
    f.data, err = FlattenFile(f.path)
  }
  return
}

func (f *FileProvider) Lookup(key string) (interface{}, error) {
  key = tryKeyfunc(key, f.Keyfunc, IdentityKey)
  d, ok := f.data[key]
  if !ok {
    return nil, nil
  }
  return d, nil
}

func FlattenMap(in map[string]interface{}) map[string]interface{} {
  f := map[string]interface{}{}
  flatten(in, "", f)
  return f
}

func FlattenFile(path string) (map[string]interface{}, error) {
  ext := filepath.Ext(path)
  switch ext {
  case ".yaml", ".yml":
    return FlattenYAMLFile(path)
  default:
    return nil, fmt.Errorf("unknown file extension: %s, expected .yaml or .yml", ext)
  }
}

func FlattenYAMLFile(path string) (map[string]interface{}, error) {
  y, err := LoadYAML(path)
  if err != nil {
    return nil, err
  }
  return FlattenMap(y), nil
}
