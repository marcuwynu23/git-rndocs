# Release Notes {{ .Version }}

{{ .ReleaseDate }}

{{ range .Sections }}{{ if .Commits }}
## {{ .Title }}

{{ range .Commits }}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Header }}{{ if .Breaking }} ⚠️{{ end }}
{{ end }}{{ end }}{{ end }}
