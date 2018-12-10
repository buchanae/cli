package cli

import (
  "encoding/json"
  "io/ioutil"
  "os"
	"github.com/ghodss/yaml"
  "github.com/BurntSushi/toml"
)

func YAMLFile(path string) *FileProvider {
  return &FileProvider{
    Path: path,
    Load: loadYaml,
  }
}

func JSONFile(path string) *FileProvider {
  return &FileProvider{
    Path: path,
    Load: loadJson,
  }
}

func TOMLFile(path string) *FileProvider {
  return &FileProvider{
    Path: path,
    Load: loadToml,
  }
}

func loadYaml(path string, data map[string]interface{}) error {
  b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
  return yaml.Unmarshal(b, &data)
}

func loadJson(path string, data map[string]interface{}) error {
  b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
  return json.Unmarshal(b, &data)
}

func loadToml(path string, data map[string]interface{}) error {
  b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
  return toml.Unmarshal(b, &data)
}

type FileProvider struct {
  KeyFunc
  Path string
  Load func(path string, data map[string]interface{}) error
  data map[string]interface{}
}

func (f *FileProvider) Init() error {
  path := os.ExpandEnv(f.Path)

  if path == "" {
    return nil
  }
  if _, err := os.Stat(path); os.IsNotExist(err) {
    return nil
  }

  data := map[string]interface{}{}
  err := f.Load(path, data)
  if err != nil {
    return err
  }

  if f.data == nil {
    f.data = map[string]interface{}{}
  }
  flatten(data, nil, f.keyfunc, f.data)
  return nil
}

func (f *FileProvider) keyfunc(key []string) string {
  if f.KeyFunc != nil {
    return f.KeyFunc(key)
  } else {
    return DotKey(key)
  }
}

func (f *FileProvider) Lookup(key []string) (interface{}, bool) {
  k := f.keyfunc(key)
  val, ok := f.data[k]
  return val, ok
}
