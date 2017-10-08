package roger

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// LoadJSON is a utility that loads a JSON file into a nested map.
func LoadJSON(path string) (map[string]interface{}, error) {
	jsonconf := map[string]interface{}{}
	jsonb, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonb, &jsonconf)
	if err != nil {
		return nil, err
	}
	return jsonconf, nil
}

// LoadYAML is a utility that loads a YAML file into a nested map.
func LoadYAML(path string) (map[string]interface{}, error) {
	yamlconf := map[string]interface{}{}
	yamlb, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlb, &yamlconf)
	if err != nil {
		return nil, err
	}
	return yamlconf, nil
}
