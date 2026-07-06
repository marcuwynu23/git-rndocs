package markdown

import (
	"strings"
	"testing"

	"github.com/marcuwynu23/git-rndocs/internal/contributors"
	"github.com/marcuwynu23/git-rndocs/internal/parser"
	"github.com/marcuwynu23/git-rndocs/internal/stats"
)

func TestBuildMarkdown(t *testing.T) {
	data := &ReleaseData{
		Version: "v1.0.0",
		Date:    "2026-07-07",
		Sections: []Section{
			{
				Title: "Features",
				Commits: []*parser.ParsedCommit{
					{Header: "add login", Type: parser.TypeFeat, Scope: "auth"},
					{Header: "add logout", Type: parser.TypeFeat},
				},
			},
			{
				Title: "Bug Fixes",
				Commits: []*parser.ParsedCommit{
					{Header: "fix crash", Type: parser.TypeFix},
				},
			},
		},
		Contributors: []contributors.Contributor{
			{Name: "Alice", Count: 3},
		},
		From: "v0.9.0",
		To:   "v1.0.0",
	}

	md := BuildMarkdown(data)

	if !strings.Contains(md, "v1.0.0") {
		t.Error("expected version in output")
	}
	if !strings.Contains(md, "Add login") {
		t.Error("expected commit message in output")
	}
	if !strings.Contains(md, "Alice") {
		t.Error("expected contributor in output")
	}
	if !strings.Contains(md, "v0.9.0...v1.0.0") {
		t.Error("expected changelog link in output")
	}
}

func TestBuildMarkdownWithStats(t *testing.T) {
	data := &ReleaseData{
		Version: "v2.0.0",
		Date:    "2026-07-07",
		Sections: []Section{
			{Title: "Features", Commits: []*parser.ParsedCommit{{Header: "new feature"}}},
		},
		Statistics: &stats.Statistics{
			TotalCommits: 1,
			ByCategory:   map[string]int{"Features": 1},
		},
	}

	md := BuildMarkdown(data)

	if !strings.Contains(md, "1 features") {
		t.Error("expected summary with 1 feature")
	}
}

func TestFormatCommit(t *testing.T) {
	c := &parser.ParsedCommit{
		Header: "add login",
		Scope:  "auth",
		Issues: []string{"123"},
		PRs:    []string{"456"},
	}

	result := formatCommit(c)
	if !strings.Contains(result, "Add login") {
		t.Error("expected header in formatted commit")
	}
}

func TestBuildMarkdownBreaking(t *testing.T) {
	data := &ReleaseData{
		Version: "v1.0.0",
		Date:    "2026-07-07",
		Breaking: []*parser.ParsedCommit{
			{Header: "API v2", Breaking: true},
		},
	}

	md := BuildMarkdown(data)

	if !strings.Contains(md, "Breaking Changes") {
		t.Error("expected Breaking Changes section")
	}
}
