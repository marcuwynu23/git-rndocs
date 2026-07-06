package template

import (
	"strings"
	"testing"

	"github.com/marcuwynu23/git-rndocs/internal/contributors"
	"github.com/marcuwynu23/git-rndocs/internal/parser"
	"github.com/marcuwynu23/git-rndocs/internal/stats"
)

func TestNewEngine(t *testing.T) {
	eng := NewEngine("")
	if eng == nil {
		t.Fatal("expected non-nil engine")
	}
}

func TestRenderDefaultTemplate(t *testing.T) {
	eng := NewEngine("")
	data := &RenderData{
		Version:     "v1.0.0",
		ReleaseDate: "2026-07-07",
		Sections: []Section{
			{
				Title: "Features",
				Commits: []*parser.ParsedCommit{
					{Header: "add login", Type: parser.TypeFeat},
				},
			},
		},
	}

	result, err := eng.Render("default", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(result, "v1.0.0") {
		t.Error("expected version in rendered output")
	}
	if !strings.Contains(result, "add login") {
		t.Error("expected commit in rendered output")
	}
}

func TestRenderMinimalTemplate(t *testing.T) {
	eng := NewEngine("")
	data := &RenderData{
		Version: "v1.0.0",
	}

	result, err := eng.Render("minimal", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(result, "v1.0.0") {
		t.Error("expected version in rendered output")
	}
}

func TestRenderGithubTemplate(t *testing.T) {
	eng := NewEngine("")
	data := &RenderData{
		Version:     "v1.0.0",
		ReleaseDate: "2026-07-07",
		Sections: []Section{
			{
				Title: "Features",
				Commits: []*parser.ParsedCommit{
					{Header: "new feature", Type: parser.TypeFeat},
				},
			},
		},
	}

	result, err := eng.Render("github", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(result, "v1.0.0") {
		t.Error("expected version in rendered output")
	}
}

func TestRenderUnknownTemplate(t *testing.T) {
	eng := NewEngine("")
	_, err := eng.Render("nonexistent", &RenderData{})

	if err == nil {
		t.Error("expected error for unknown template")
	}
}

func TestGetBuiltin(t *testing.T) {
	if GetBuiltin("default") == "" {
		t.Error("expected non-empty default template")
	}
	if GetBuiltin("github") == "" {
		t.Error("expected non-empty github template")
	}
	if GetBuiltin("minimal") == "" {
		t.Error("expected non-empty minimal template")
	}
	if GetBuiltin("unknown") != "" {
		t.Error("expected empty for unknown template")
	}
}

func TestBuildRenderData(t *testing.T) {
	eng := NewEngine("")
	categorized := map[parser.CommitType][]*parser.ParsedCommit{
		parser.TypeFeat: {{Header: "feature", Type: parser.TypeFeat}},
	}
	contribs := []contributors.Contributor{{Name: "Alice", Count: 1}}
	stat := &stats.Statistics{TotalCommits: 1}

	data := eng.BuildRenderData("v1.0.0", "2026-07-07", "summary", "v0.9.0", "v1.0.0", categorized, nil, contribs, stat)

	if data.Version != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", data.Version)
	}
	if len(data.Sections) != 1 {
		t.Errorf("expected 1 section, got %d", len(data.Sections))
	}
}
