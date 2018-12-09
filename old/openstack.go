package roger

import (
  "strings"
  "net/http"
)

type OpenstackMetadataProvider struct {
  URL string
  vals map[string]interface{}
}

func Openstack() *OpenstackMetadataProvider {
  return &OpenstackMetadataProvider{
    URL: "http://169.254.169.254/openstack/latest/meta_data.json",
  }
}

func (g *OpenstackMetadataProvider) Init() error {
  g.vals = map[string]interface{}{}

  // TODO this should have a quick timeout
  req, err := http.NewRequest("GET", g.URL, nil)
  if err != nil {
    return err
  }

  m, err := curlJSON(req, "openstack")
  if err != nil {
    return err
  }

  for k, v := range m {
    k = NormalizeKey(k)
    k = strings.TrimPrefix(k, "openstack.meta.")
    g.vals[k] = v
  }
  return nil
}

func (g *OpenstackMetadataProvider) Lookup(key string) (interface{}, error) {
  key = NormalizeKey(key)
  d, ok := g.vals[key]
  if !ok {
    return nil, nil
  }
  return d, nil
}
