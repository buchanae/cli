package roger

import (
	"fmt"
	"os"
	"path/filepath"
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
func (f *FileProvider) Init() (err error) {
	if f.path != "" {
		f.data, err = FlattenFile(f.path)
	}
	return
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

// FlattenMap is a utility used to flatten a nested map. For example:
//   "root": {
//     "sub": {
//       "subone": "val",
//     },
//   }
//
// flattens to: {"root.sub.subone": "val"}
func FlattenMap(in map[string]interface{}) map[string]interface{} {
	f := map[string]interface{}{}
	flatten(in, "", f)
	return f
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

// FlattenFile is a utility that loads a file into a nested map and calles FlattenMap.
// The file type is determined by the extension. Currently only YAML is supported,
// with ".yaml" and ".yml" extension.
func FlattenFile(path string) (map[string]interface{}, error) {
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		return FlattenYAMLFile(path)
	default:
		return nil, fmt.Errorf("unknown file extension: %s, expected .yaml or .yml", ext)
	}
}

// FlattenYAMLFile is a utility that loads a YAML file into a nested map and calls FlattenMap.
func FlattenYAMLFile(path string) (map[string]interface{}, error) {
	y, err := LoadYAML(path)
	if err != nil {
		return nil, err
	}
	return FlattenMap(y), nil
}
