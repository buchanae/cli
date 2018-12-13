package cli

import (
  "os"
  "io/ioutil"
  "github.com/ghodss/yaml"
  "encoding/json"
  "github.com/BurntSushi/toml"
)

var DefaultYAML = FileOpts{
  Paths: []string{"config.yaml", "config.yml"},
  OptKey: []string{"config"},
}

var DefaultJSON = FileOpts{
  Paths: []string{"config.json"},
  OptKey: []string{"config"},
}

var DefaultTOML = FileOpts{
  Paths: []string{"config.toml"},
  OptKey: []string{"config"},
}


type FileOpts struct {
  Paths []string
  OptKey []string
}


func YAML(opts FileOpts) Provider {
  return &fileProvider{opts, yaml.Unmarshal}
}

func JSON(opts FileOpts) Provider {
  return &fileProvider{opts, json.Unmarshal}
}

func TOML(opts FileOpts) Provider {
  return &fileProvider{opts, toml.Unmarshal}
}

type fileProvider struct {
  opts FileOpts
  unm unmarshaler
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
