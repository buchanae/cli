package cli

import (
  "os"
  "io/ioutil"
  "github.com/ghodss/yaml"
  "encoding/json"
  "github.com/BurntSushi/toml"
)

var DefaultYAML = &YAML{
  Paths: []string{"config.yaml", "config.yml"},
  OptKey: []string{"config"},
}

var DefaultJSON = &JSON{
  Paths: []string{"config.json"},
  OptKey: []string{"config"},
}

var DefaultTOML = &TOML{
  Paths: []string{"config.toml"},
  OptKey: []string{"config"},
}


type YAML struct {
  Paths []string
  OptKey []string
}

type JSON struct {
  Paths []string
  OptKey []string
}

type TOML struct {
  Paths []string
  OptKey []string
}


func (y *YAML) Load(l *Loader) error {
  return loadFile(y.OptKey, y.Paths, yaml.Unmarshal, l)
}

func (j *JSON) Load(l *Loader) error {
  return loadFile(j.OptKey, j.Paths, json.Unmarshal, l)
}

func (t *TOML) Load(l *Loader) error {
  return loadFile(t.OptKey, t.Paths, toml.Unmarshal, l)
}


type unmarshaler func([]byte, interface{}) error

func loadFile(optKey []string, paths []string, unm unmarshaler, l *Loader) error {

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
    err = unm(b, &data)
    if err != nil {
      return err
    }

	  flatten2(data, l, nil)
  }

  return nil
}
