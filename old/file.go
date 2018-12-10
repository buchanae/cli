package roger

import (
	"fmt"
  "io/ioutil"
	"os"
)

// FileProvider provides access to values configured via a file.
// Currently only YAML is supported.
type FileProvider struct {
	Keyfunc
	Paths []string
	data map[string]interface{}
}

func Files(paths ...string) *FileProvider {
	return &FileProvider{Paths: paths}
}

func loadFile(path string, u unmarshaler, flat map[string]interface{}) error {
  conf := map[string]interface{}{}

  err = u(b, &conf)
	if err != nil {
		return err
	}
  flatten(conf, "", flat)
  return nil
}
