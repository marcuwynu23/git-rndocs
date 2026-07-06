package stats

import (
	"context"
	"testing"

	"github.com/marcuwynu23/git-rndocs/internal/parser"
)

func TestCollectBasic(t *testing.T) {
	categorized := map[parser.CommitType][]*parser.ParsedCommit{
		parser.TypeFeat: {{Header: "f1"}, {Header: "f2"}},
		parser.TypeFix:  {{Header: "x1"}},
	}

	stats, err := Collect(context.Background(), nil, "", "", categorized)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if stats.TotalCommits != 3 {
		t.Errorf("expected 3 total commits, got %d", stats.TotalCommits)
	}
	if stats.ByCategory["Features"] != 2 {
		t.Errorf("expected 2 features, got %d", stats.ByCategory["Features"])
	}
	if stats.ByCategory["Bug Fixes"] != 1 {
		t.Errorf("expected 1 bug fix, got %d", stats.ByCategory["Bug Fixes"])
	}
}

func TestCollectEmpty(t *testing.T) {
	stats, err := Collect(context.Background(), nil, "", "", nil)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if stats.TotalCommits != 0 {
		t.Errorf("expected 0 commits, got %d", stats.TotalCommits)
	}
}
