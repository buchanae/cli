package main

var rawtpl = `
package {{ .Pkgname }}

import "github.com/buchanae/roger/roger"

func (c *Config) RogerVals() roger.Vals {
  m := map[string]roger.Val{
    {{ range .Nodes -}}
      "{{ join .Key }}": roger.NewVal("{{ synopsis .Doc }}", &c.{{ join .Key }}),
    {{ end }}
  }
  {{ range $index, $element := .Alias -}}
    if v, ok := m["{{ $element }}"]; ok {
      m["{{ $index }}"] = v
    }
  {{ end }}
  return m
}
`
