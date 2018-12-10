package inspect

import (
  "text/template"
)

var DefaultTemplate = template.Must(template.New("default").Parse(`
package {{ .Package }}

{{ range $name, $path := .Imports -}}
import {{ $name }} "{{ $path }}"
{{ end }}

var cmdSpecs = []cli.CmdSpec{
{{ range .Funcs -}}
  &{{ .FuncNamePriv }}Spec{
    {{ if .HasDefaultOpts -}}
    Opt: {{ .DefaultOptsName }},
    {{- end }}
  },
{{ end }}
}

{{ range .Funcs }}
type {{ .FuncNamePriv }}Spec struct {
  {{- if .HasOpts }}
  Opt {{ .OptsType }}
  {{ end }}

  args struct {
    {{ range .Args }}
      arg{{ .Idx }} {{ .Type }}
    {{ end }}
  }
}

func (cmd *{{ .FuncNamePriv }}Spec) Name() string {
  return "{{ .FuncName }}"
}

func (cmd *{{ .FuncNamePriv }}Spec) Doc() string {
  return {{ .Doc | printf "%q" }}
}

func (cmd *{{ .FuncNamePriv }}Spec) Run(args []string) {
  cli.CheckArgs(args, cmd.ArgSpecs())
  {{ .FuncName }}(
  {{- if .HasOpts }}
    cmd.Opt,
  {{ end -}}
  {{- range .Args -}}
    {{ if .Variadic -}}
    cmd.args.arg{{ .Idx }}...,
    {{- else -}}
    cmd.args.arg{{ .Idx }},
    {{- end }}
  {{ end -}}
  )
}

func (cmd *{{ .FuncNamePriv }}Spec) ArgSpecs() []cli.ArgSpec {
  {{ if not .HasArgs }}
  return nil
  {{ else -}}
  return []cli.ArgSpec{
    {{ range .Args -}}
    {
      Name: "{{ .Name }}",
      Type: "{{ .Type }}",
      Variadic: {{ .Variadic }},
      Value: &cmd.args.arg{{ .Idx }},
    },
    {{- end }}
  }
  {{- end }}
}

func (cmd *{{ .FuncNamePriv }}Spec) OptSpecs() []cli.OptSpec {
  {{ if not .HasOpts }}
  return nil
  {{ else -}}
  return []cli.OptSpec{
    {{ range .Opts -}}
    {
      Key: {{ .Key | printf "%#v" }},
      Doc: {{ .Doc | printf "%q" }},
      Value: &cmd.Opt.{{ .KeyJoined }},
    },
    {{- end }}
  }
  {{- end }}
}
{{ end }}
`))
