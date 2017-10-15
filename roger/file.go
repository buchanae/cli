package roger

import (
	"fmt"
  "io/ioutil"
	"os"
	"path/filepath"
	"encoding/json"
	"github.com/ghodss/yaml"
  "github.com/BurntSushi/toml"
)

// FileProvider provides access to values configured via a file.
// Currently only YAML is supported.
type FileProvider struct {
	Keyfunc
	path string
	data map[string]interface{}
}

// NewFileProvider returns a FileProvider for the given path.
// The type of the file is determined by the file extension.
// Currently only YAML is supported, via ".yaml" and ".yml".
func NewFileProvider(path string) *FileProvider {
	return &FileProvider{path: path}
}

// OptionalFileProvider returns a FileProvider, where the file is optional.
// If path is "" or is missing, the provider will do nothing.
func OptionalFileProvider(path string) *FileProvider {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = ""
	}
	return NewFileProvider(path)
}

// Init loads the file.
func (f *FileProvider) Init() error {
  f.data = map[string]interface{}{}

	if f.path != "" {
    ext := filepath.Ext(f.path)
    switch ext {
    case ".yaml", ".yml":
      return loadFile(f.path, yaml.Unmarshal, f.data)
    case ".json":
      return loadFile(f.path, json.Unmarshal, f.data)
    case ".toml":
      return loadFile(f.path, toml.Unmarshal, f.data)
    default:
      return fmt.Errorf("unknown file extension: %s, expected .yaml or .yml", ext)
    }
	}
	return nil
}

// Lookup looks up values in the file.
func (f *FileProvider) Lookup(key string) (interface{}, error) {
	key = tryKeyfunc(key, f.Keyfunc, IdentityKey)
	d, ok := f.data[key]
	if !ok {
		return nil, nil
	}
	return d, nil
}

// flatten flattens a nested map. For example:
//   "root": {
//     "sub": {
//       "subone": "val",
//     },
//   }
//
// flattens to: {"root.sub.subone": "val"}
func flatten(in map[string]interface{}, prefix string, out map[string]interface{}) {
	for k, v := range in {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}

		switch x := v.(type) {
		case map[string]interface{}:
			flatten(x, path, out)
		default:
			out[path] = v
		}
	}
}

func loadFile(path string, u unmarshaler, flat map[string]interface{}) error {
  conf := map[string]interface{}{}
  b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

  err = u(b, &conf)
	if err != nil {
		return err
	}
  flatten(conf, "", flat)
  return nil
}

type unmarshaler func([]byte, interface{}) error
