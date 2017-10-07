package roger

import (
  "fmt"
  "path/filepath"
)

type FileProvider struct {
  Keyfunc
  data map[string]interface{}
}

func NewFileProvider(path string) (*FileProvider, error) {
  data, err := FlattenFile(path)
  if err != nil {
    return nil, err
  }
  return &FileProvider{data: data}, nil
}

func (f *FileProvider) Lookup(key string) (interface{}, bool) {
  key = tryKeyfunc(key, f.Keyfunc, IdentityKey)
  d, ok := f.data[key]
  return d, ok
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
    return nil, fmt.Errorf("unknown file type: %s", ext)
  }
}

func FlattenYAMLFile(path string) (map[string]interface{}, error) {
  y, err := LoadYAML(path)
  if err != nil {
    return nil, err
  }
  return FlattenMap(y), nil
}
