package cli

import (
  "fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
)

// YAMLFile creates a provider that loads values from a YAML file.
func YAMLFile(path string) Provider {
	return &fileProvider{
		Path: path,
		Load: loadYaml,
	}
}

// JSONFile creates a provider that loads values from a JSON file.
func JSONFile(path string) Provider {
	return &fileProvider{
		Path: path,
		Load: loadJson,
	}
}

// TOMLFile creates a provider that loads values from a TOML file.
func TOMLFile(path string) Provider {
	return &fileProvider{
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

// fileProvider provides option values from a file.
type fileProvider struct {
	KeyFunc
	Path string
	Load func(path string, data map[string]interface{}) error
	data map[string]interface{}
}

func (f *fileProvider) Init() error {
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
	flatten(data, f.data, nil, f.keyfunc)
	return nil
}

func (f *fileProvider) ValidateKeys(keys [][]string) []error {
  var errs []error

  allowed := map[string]struct{}{}
  for _, key := range keys {
    k := f.keyfunc(key)
    allowed[k] = struct{}{}
  }

  for k, _ := range f.data {
    _, ok := allowed[k]
    if !ok {
      errs = append(errs, fmt.Errorf("unrecognized option %q in file %q", k, f.Path))
    }
  }

  return errs
}

func (f *fileProvider) keyfunc(key []string) string {
	if f.KeyFunc != nil {
		return f.KeyFunc(key)
	} else {
		return DotKey(key)
	}
}

func (f *fileProvider) Lookup(key []string) (interface{}, bool) {
	k := f.keyfunc(key)
	val, ok := f.data[k]
	return val, ok
}
