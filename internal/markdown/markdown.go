package markdown

import (
	"fmt"
	"strings"

	"github.com/marcuwynu23/git-rndocs/internal/parser"
	"github.com/marcuwynu23/git-rndocs/internal/stats"
	"github.com/marcuwynu23/git-rndocs/internal/contributors"
)

type ReleaseData struct {
	Version      string
	Date         string
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

func BuildMarkdown(data *ReleaseData) string {
	var b strings.Builder

	b.WriteString("# Release Notes\n\n")
	b.WriteString(fmt.Sprintf("## Version\n\n%s\n\n", data.Version))
	b.WriteString(fmt.Sprintf("## Release Date\n\n%s\n\n", data.Date))

	if data.Statistics != nil && data.Summary == "" {
		var parts []string
		for _, s := range data.Sections {
			if len(s.Commits) > 0 {
				parts = append(parts, fmt.Sprintf("- %d %s", len(s.Commits), strings.ToLower(s.Title)))
			}
		}
		if len(parts) > 0 {
			b.WriteString("## Summary\n\n")
			b.WriteString(strings.Join(parts, "\n"))
			b.WriteString("\n\n")
		}
	}

	if data.Summary != "" {
		b.WriteString("## Summary\n\n")
		b.WriteString(data.Summary)
		b.WriteString("\n\n")
	}

	b.WriteString("---\n\n")

	allBreaking := data.Breaking
	for _, sec := range data.Sections {
		sectionCommits := sec.Commits
		if len(sectionCommits) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("## %s\n\n", sec.Title))
		for _, c := range sectionCommits {
			b.WriteString(formatCommit(c))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(allBreaking) > 0 {
		b.WriteString("## Breaking Changes\n\n")
		for _, c := range allBreaking {
			b.WriteString(formatCommit(c))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(data.Contributors) > 0 {
		b.WriteString("## Contributors\n\n")
		for _, c := range data.Contributors {
			b.WriteString(fmt.Sprintf("- %s (%d commits)\n", c.Name, c.Count))
		}
		b.WriteString("\n")
	}

	if data.From != "" && data.To != "" {
		b.WriteString("## Full Changelog\n\n")
		b.WriteString(fmt.Sprintf("%s...%s\n", data.From, data.To))
		b.WriteString("\n")
	}

	return b.String()
}

func formatCommit(c *parser.ParsedCommit) string {
	prefix := "- "
	scope := ""
	if c.Scope != "" {
		scope = fmt.Sprintf(" **%s:**", c.Scope)
	}
	header := c.Header
	if header == "" {
		header = c.Raw
	}
	if len(header) > 0 {
		header = strings.ToUpper(header[:1]) + header[1:]
	}
	msg := fmt.Sprintf("%s%s%s", prefix, scope, header)

	var refs []string
	for _, issue := range c.Issues {
		refs = append(refs, fmt.Sprintf("#%s", issue))
	}
	for _, pr := range c.PRs {
		refs = append(refs, fmt.Sprintf("#%s", pr))
	}
	if len(refs) > 0 {
		msg += fmt.Sprintf(" (%s)", strings.Join(refs, ", "))
	}

	if c.Breaking {
		msg += " [breaking]"
	}

	return msg
}

type TemplateData struct {
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
