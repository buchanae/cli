package cli

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

// DefaultYAML contains the most common configuration
// for loading options from a YAML config file.
var DefaultYAML = FileOpts{
	Paths:  []string{"config.yaml", "config.yml"},
	OptKey: []string{"config"},
}

// DefaultJSON contains the most common configuration
// for loading options from a JSON config file.
var DefaultJSON = FileOpts{
	Paths:  []string{"config.json"},
	OptKey: []string{"config"},
}

// DefaultTOML contains the most common configuration
// for loading options from a TOML config file.
var DefaultTOML = FileOpts{
	Paths:  []string{"config.toml"},
	OptKey: []string{"config"},
}

// FileOpts describes options related to loading
// options from a file.
type FileOpts struct {
	// Paths is a list of paths to look for a config file.
	// Environment variables will be expanded using `os.ExpandEnv`.
	// Loading will stop on the first path that exists.
	Paths []string
	// OptKey is used to look for a config file path set by
	// an option (e.g. by a flag or env. var). For example,
	// an OptKey of ["config", "file"] could load the config
	// file from a "--config.file" flag. OptKey is prepended
	// to the Paths list, so it takes priority.
	OptKey []string
}

// YAML loads options from a YAML file.
func YAML(opts FileOpts) Provider {
	return &fileProvider{opts, yaml.Unmarshal}
}

// JSON loads options from a JSON file.
func JSON(opts FileOpts) Provider {
	return &fileProvider{opts, json.Unmarshal}
}

// TOML loads options from a TOML file.
func TOML(opts FileOpts) Provider {
	return &fileProvider{opts, toml.Unmarshal}
}

type fileProvider struct {
	opts FileOpts
	unm  unmarshaler
}

func (f *fileProvider) Provide(l *Loader) error {
	optKey := f.opts.OptKey
	paths := f.opts.Paths
	opt := l.GetString(optKey)
	paths = append([]string{opt}, paths...)

	for _, path := range paths {
		path := os.ExpandEnv(path)
		if path == "" || !exists(path) {
			continue
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		data := map[string]interface{}{}
		err = f.unm(b, &data)
		if err != nil {
			return err
		}

		flatten2(data, l, nil)
	}

	return nil
}

type unmarshaler func([]byte, interface{}) error
