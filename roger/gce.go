package roger

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
  "net/http"
  "strings"
)

type GCEMetadataProvider struct {
  URL string
  vals map[string]interface{}
}

func NewGCEMetadataProvider() *GCEMetadataProvider {
  return &GCEMetadataProvider{URL: "http://metadata.google.internal"}
}

func (g *GCEMetadataProvider) Init() error {
  g.vals = map[string]interface{}{}
  // TODO this should have a quick timeout.
  m, err := LoadGCEMetadata(g.URL)
  if err != nil {
    return err
  }
  f := map[string]interface{}{}
  flatten(m, "gce", f)

  for k, v := range f {
    k = NormalizeKey(k)
    k = strings.TrimPrefix(k, "gce.instance.attributes.")
    fmt.Println("GCE", k, v)
    g.vals[k] = v
  }
  return nil
}

func (g *GCEMetadataProvider) Lookup(key string) (interface{}, error) {
  key = NormalizeKey(key)
  d, ok := g.vals[key]
  if !ok {
    return nil, nil
  }
  return d, nil
}

// LoadGCEMetadata loads metadata from the given URL.
func LoadGCEMetadata(url string) (map[string]interface{}, error) {
  meta := map[string]interface{}{}
	client := http.Client{}
	path := "/computeMetadata/v1/?recursive=true"
	req, err := http.NewRequest("GET", url+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 response status from GCE Metadata: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	perr := json.Unmarshal(body, &meta)
	if perr != nil {
		return nil, perr
	}
	return meta, nil
}
