package main

var rawtpl = `
package {{ .Pkgname }}

import "github.com/buchanae/roger/roger"

// TODO how to determine if this should have a pointer receiver?
func (c *Config) RogerVals() map[string]roger.Val {
  return map[string]roger.Val{
    {{ range .Nodes -}}
      "{{ join .Key }}": roger.NewVal({{ keysyn .Key }}, "{{ synopsis .Doc }}", &c.{{ join .Key }}),
    {{ end }}
  }
}
`
