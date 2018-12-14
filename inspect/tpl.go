package inspect

import (
	"text/template"
)

var DefaultTemplate = template.Must(template.New("default").Parse(`
package {{ .Package }}

{{ range $name, $path := .Imports -}}
import {{ $name }} "{{ $path }}"
{{ end }}

func specs() []cli.Spec {
  return []cli.Spec{
  {{ range .Funcs -}}
    &{{ .FuncNamePriv }}Spec{
      {{ if .HasDefaultOpts -}}
      opt: {{ .DefaultOptsName }},
      {{- end }}
    },
  {{ end }}
  }
}

{{ range .Funcs }}
type {{ .FuncNamePriv }}Spec struct {
  cmd *cli.Cmd
  {{ if .HasOpts -}}
  opt {{ .OptsType }}
  {{- end }}
  args struct {
    {{ range .Args -}}
      arg{{ .Idx }} {{ .Type }}
    {{ end }}
  }
}

func (cmd *{{ .FuncNamePriv }}Spec) Run() {
  {{ .FuncName }}(
  {{- if .HasOpts }}
    cmd.opt,
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

func (cmd *{{ .FuncNamePriv }}Spec) Cmd() *cli.Cmd {
  if cmd.cmd != nil {
    return cmd.cmd
  }
  cmd.cmd = &cli.Cmd{
    FuncName:   {{ .FuncName | printf "%q" }},
    RawDoc: {{ .Doc | printf "%q" }},
    Args: []*cli.Arg{
      {{ range .Args -}}
      {
        Name: "{{ .Name }}",
        Type: "{{ .Type }}",
        Variadic: {{ .Variadic }},
        Value: &cmd.args.arg{{ .Idx }},
      },
      {{- end }}
    },
    Opts: []*cli.Opt{
      {{ range .Opts -}}
      {
        Key: {{ .Key | printf "%#v" }},
        RawDoc: {{ .Doc | printf "%q" }},
        Value: &cmd.opt.{{ .KeyJoined }},
        DefaultValue: cmd.opt.{{ .KeyJoined }},
        Type: {{ .Type | printf "%q" }},
        Short: {{ .Short | printf "%q" }},
      },
      {{- end }}
    },
  }
  return cmd.cmd
}
{{ end }}
`))
