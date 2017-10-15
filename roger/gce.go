package roger

import (
  "strings"
  "net/http"
)

type GCEMetadataProvider struct {
  URL string
  vals map[string]interface{}
}

func GCE() *GCEMetadataProvider {
  return &GCEMetadataProvider{
    URL: "http://metadata.google.internal/computeMetadata/v1/?recusive=true",
  }
}

func (g *GCEMetadataProvider) Init() error {
  g.vals = map[string]interface{}{}

  // TODO this should have a quick timeout
  req, err := http.NewRequest("GET", g.URL, nil)
	req.Header.Add("Metadata-Flavor", "Google")
  if err != nil {
    return err
  }

  m, err := curlJSON(req, "gce")
  if err != nil {
    return err
  }

  for k, v := range m {
    k = NormalizeKey(k)
    k = strings.TrimPrefix(k, "gce.instance.attributes.")
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
