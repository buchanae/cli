package {{ .Package }}

{{ range $name, $path := .Imports -}}
import {{ $name }} "{{ $path }}"
{{ end }}

{{ range .Funcs }}
func new{{ .FuncName }}Spec() *{{ .FuncNamePriv }}Spec {
  return &{{ .FuncNamePriv }}Spec{
    CmdSpec: cli.CmdSpec{
      Name: "{{ .FuncName }}",
      Doc: {{ .Doc | printf "%q" }},
    },
    {{ if .HasDefaultOpts -}}
    Opt: {{ .DefaultOptsName }},
    {{- end }}
  }
}

type {{ .FuncNamePriv }}Spec struct {
  cli.CmdSpec
  {{- if .HasOpts }}
  Opt {{ .OptsType }}
  {{ end }}
}

func (cmd *{{ .FuncNamePriv }}Spec) Run(args []string) {
  cli.CheckArgs(args, cmd.ArgSpecs())
  {{ .FuncName }}(
  {{- if .HasOpts }}
    cmd.Opt,
  {{ end -}}
  {{- range .Args -}}
    {{ if .Variadic -}}
    cli.Coerce{{ .CoerceType }}(args[{{ .Idx }}:])...,
    {{- else -}}
    cli.Coerce{{ .CoerceType }}(args[{{ .Idx }}]),
    {{- end }}
  {{ end -}}
  )
}

func (cmd *{{ .FuncNamePriv }}Spec) ArgSpecs() []cli.ArgSpec {
  {{ if not .HasArgs }}
  return nil
  {{ else }}
  return []cli.ArgSpec{
    {{ range .Args -}}
    {
      Name: "{{ .Name }}",
      Type: "{{ .Type }}",
      Variadic: {{ .Variadic }},
    },
    {{- end }}
  }
  {{ end }}
}

func (cmd *{{ .FuncNamePriv }}Spec) OptSpecs() []cli.OptSpec {
  {{ if not .HasOpts }}
  return nil
  {{ else }}
  return []cli.OptSpec{
    {{ range .Opts -}}
    {
      Key: {{ .Key | printf "%#v" }},
      Doc: {{ .Doc | printf "%q" }},
      Type: "{{ .Type }}",
      Value: &cmd.Opt.{{ .KeyJoined }},
    },
    {{- end }}
  }
  {{ end }}
}
{{ end }}
