# {{ .Version }}

{{ .ReleaseDate }}

{{ range .Sections }}{{ if .Commits }}
{{ .Title }}:
{{ range .Commits }}
- {{ .Header }}
{{ end }}{{ end }}{{ end }}
