package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/marcuwynu23/git-rndocs/internal/contributors"
	"github.com/marcuwynu23/git-rndocs/internal/parser"
	"github.com/marcuwynu23/git-rndocs/internal/stats"
)

type Engine struct {
	customDir string
}

func NewEngine(customDir string) *Engine {
	return &Engine{customDir: customDir}
}

func (e *Engine) Render(name string, data *RenderData) (string, error) {
	tmplContent, err := e.loadTemplate(name)
	if err != nil {
		return "", err
	}

	funcMap := template.FuncMap{
		"now":     time.Now,
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"title":   strings.Title,
		"count":   func(items interface{}) int { return 0 },
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

func (e *Engine) loadTemplate(name string) (string, error) {
	if e.customDir != "" {
		customPath := filepath.Join(e.customDir, name+".md")
		if data, err := os.ReadFile(customPath); err == nil {
			return string(data), nil
		}
	}

	embedded := GetBuiltin(name)
	if embedded != "" {
		return embedded, nil
	}

	return "", fmt.Errorf("template %s not found", name)
}

type RenderData struct {
	Version      string
	ReleaseDate  string
	Summary      string
	Sections     []Section
	Breaking     []*parser.ParsedCommit
	Contributors []contributors.Contributor
	Statistics   *stats.Statistics
	From         string
	To           string
}

type Section struct {
	Title   string
	Commits []*parser.ParsedCommit
}

func (e *Engine) BuildRenderData(version, date, summary, from, to string, categorized map[parser.CommitType][]*parser.ParsedCommit, breaking []*parser.ParsedCommit, contribs []contributors.Contributor, stat *stats.Statistics) *RenderData {
	var sections []Section
	for _, ct := range []parser.CommitType{
		parser.TypeFeat, parser.TypeFix, parser.TypePerf, parser.TypeDocs,
		parser.TypeRefactor, parser.TypeStyle, parser.TypeBuild, parser.TypeCI,
		parser.TypeTest, parser.TypeChore, parser.TypeRevert, parser.TypeUnknown,
	} {
		if commits, ok := categorized[ct]; ok && len(commits) > 0 {
			sections = append(sections, Section{
				Title:   ct.DisplayName(),
				Commits: commits,
			})
		}
	}

	return &RenderData{
		Version:      version,
		ReleaseDate:  date,
		Summary:      summary,
		Sections:     sections,
		Breaking:     breaking,
		Contributors: contribs,
		Statistics:   stat,
		From:         from,
		To:           to,
	}
}

func GetBuiltin(name string) string {
	switch name {
	case "default":
		return defaultTemplate
	case "github":
		return githubTemplate
	case "minimal":
		return minimalTemplate
	default:
		return ""
	}
}

const defaultTemplate = `# Release Notes

## Version

{{ .Version }}

## Release Date

{{ .ReleaseDate }}

{{ if .Summary }}
## Summary

{{ .Summary }}

{{ end }}
---

{{ range .Sections }}
{{ if .Commits }}
## {{ .Title }}

{{ range .Commits }}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Header }}{{ if .Breaking }} [breaking]{{ end }}{{ if .Issues }} ({{ range $i, $v := .Issues }}#{{ $v }}{{ end }}){{ end }}
{{ end }}

{{ end }}
{{ end }}

{{ if .Breaking }}
## Breaking Changes

{{ range .Breaking }}
- {{ .Header }}
{{ end }}

{{ end }}

{{ if .Contributors }}
## Contributors

{{ range .Contributors }}
- {{ .Name }} ({{ .Count }} commits)
{{ end }}

{{ end }}

{{ if .From }}
## Full Changelog

{{ .From }}...{{ .To }}
{{ end }}
`

const githubTemplate = `# Release Notes {{ .Version }}

{{ .ReleaseDate }}

{{ if .Summary }}
{{ .Summary }}
{{ end }}

{{ range .Sections }}
{{ if .Commits }}
## {{ .Title }}

{{ range .Commits }}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Header }}{{ if .Breaking }} [breaking]{{ end }} {{ if .Issues }}({{ range $i, $v := .Issues }}#{{ $v }}{{ end }}){{ end }}
{{ end }}

{{ end }}
{{ end }}

{{ if .Breaking }}
## ⚠️ Breaking Changes

{{ range .Breaking }}
- {{ .Header }}
{{ end }}

{{ end }}

{{ if .Contributors }}
## Contributors

{{ range .Contributors }}
- {{ .Name }}
{{ end }}

{{ end }}
`

const minimalTemplate = `# {{ .Version }}

{{ .ReleaseDate }}

{{ range .Sections }}
{{ if .Commits }}
## {{ .Title }}

{{ range .Commits }}
- {{ .Header }}
{{ end }}

{{ end }}
{{ end }}
`
