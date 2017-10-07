package main

var rawtpl = `
package {{ .Pkgname }}

import "github.com/buchanae/roger/roger"

// TODO how to determine if this should have a pointer receiver?
func (c *Config) RogerVals() roger.Vals {
  return map[string]roger.Val{
    {{ range .Nodes -}}
      "{{ join .Key }}": roger.NewVal("{{ synopsis .Doc }}", &c.{{ join .Key }}),
    {{ end }}
  }
}
`
