package roger

import (
  "fmt"
	"encoding/json"
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
