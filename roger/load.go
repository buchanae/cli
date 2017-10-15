package roger

import (
  "fmt"
	"encoding/json"
	"github.com/ghodss/yaml"
  "net/http"
	"io/ioutil"
  "time"
)

func curlJSON(req *http.Request, prefix string) (map[string]interface{}, error) {
	client := http.Client{
    Timeout: time.Second,
  }
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 response status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

  meta := map[string]interface{}{}
	perr := json.Unmarshal(body, &meta)
	if perr != nil {
		return nil, perr
	}

  f := map[string]interface{}{}
  flatten(meta, prefix, f)
	return f, nil
}

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
