package parser

import (
	"testing"
)

func TestParseConventionalCommit(t *testing.T) {
	msg := "feat(auth): add login endpoint\n\nThis implements OAuth2 login\n\nBREAKING CHANGE: API v1 deprecated\n\nCloses #123"
	pc := ParseCommit(msg, "abc123", "Alice", "alice@example.com")

	if pc.Type != TypeFeat {
		t.Errorf("expected TypeFeat, got %s", pc.Type)
	}
	if pc.Scope != "auth" {
		t.Errorf("expected scope 'auth', got '%s'", pc.Scope)
	}
	if !pc.Breaking {
		t.Error("expected breaking change")
	}
	if pc.Header != "add login endpoint" {
		t.Errorf("expected header 'add login endpoint', got '%s'", pc.Header)
	}
	if len(pc.Issues) == 0 || pc.Issues[0] != "123" {
		t.Errorf("expected issue #123, got %v", pc.Issues)
	}
}

func TestParseNonConventionalCommit(t *testing.T) {
	msg := "fixed the thing"
	pc := ParseCommit(msg, "def456", "Bob", "bob@test.com")

	if pc.Type != TypeUnknown {
		t.Errorf("expected TypeUnknown, got %s", pc.Type)
	}
}

func TestParseBreakingWithBang(t *testing.T) {
	msg := "feat!: new API"
	pc := ParseCommit(msg, "", "", "")

	if !pc.Breaking {
		t.Error("expected breaking change with !")
	}
}

func TestParseScope(t *testing.T) {
	msg := "fix(core): resolve null pointer"
	pc := ParseCommit(msg, "", "", "")

	if pc.Scope != "core" {
		t.Errorf("expected scope 'core', got '%s'", pc.Scope)
	}
	if pc.Type != TypeFix {
		t.Errorf("expected TypeFix, got %s", pc.Type)
	}
}

func TestCategorizeCommits(t *testing.T) {
	commits := []*ParsedCommit{
		{Type: TypeFeat, Header: "feature 1"},
		{Type: TypeFix, Header: "fix 1"},
		{Type: TypeFeat, Header: "feature 2"},
		{Type: TypeChore, Header: "chore 1"},
	}

	include := []string{"feat", "fix"}
	exclude := []string{}

	categorized := CategorizeCommits(commits, include, exclude)

	if len(categorized[TypeFeat]) != 2 {
		t.Errorf("expected 2 feat commits, got %d", len(categorized[TypeFeat]))
	}
	if len(categorized[TypeFix]) != 1 {
		t.Errorf("expected 1 fix commit, got %d", len(categorized[TypeFix]))
	}
	if _, ok := categorized[TypeChore]; ok {
		t.Error("expected chore to be excluded")
	}
}

func TestDisplayName(t *testing.T) {
	tests := []struct {
		ct   CommitType
		want string
	}{
		{TypeFeat, "Features"},
		{TypeFix, "Bug Fixes"},
		{TypeDocs, "Documentation"},
		{TypeUnknown, "Other Changes"},
	}

	for _, tt := range tests {
		if got := tt.ct.DisplayName(); got != tt.want {
			t.Errorf("DisplayName(%s) = %s, want %s", tt.ct, got, tt.want)
		}
	}
}

func TestExtractPRs(t *testing.T) {
	msg := "feat: add feature (#42)"
	pc := ParseCommit(msg, "", "", "")

	if len(pc.PRs) == 0 || pc.PRs[0] != "42" {
		t.Errorf("expected PR #42, got %v", pc.PRs)
	}
}

func TestParseFooters(t *testing.T) {
	lines := []string{
		"feat: add feature",
		"Reviewed-by: Alice",
		"Refs: #123",
	}
	footers := parseFooters(lines)

	if len(footers) != 2 {
		t.Errorf("expected 2 footers, got %d", len(footers))
	}
	if footers[0].Token != "Reviewed-by" {
		t.Errorf("expected Reviewed-by, got %s", footers[0].Token)
	}
}
