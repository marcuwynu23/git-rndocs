package contributors

import (
	"context"
	"testing"

	"github.com/marcuwynu23/git-rndocs/internal/parser"
)

func TestCollect(t *testing.T) {
	commits := []*parser.ParsedCommit{
		{Author: "Alice", AuthorEmail: "alice@test.com"},
		{Author: "Bob", AuthorEmail: "bob@test.com"},
		{Author: "Alice", AuthorEmail: "alice@test.com"},
	}

	contribs := Collect(context.Background(), nil, "", "", commits)

	if len(contribs) != 2 {
		t.Errorf("expected 2 contributors, got %d", len(contribs))
	}

	if contribs[0].Name != "Alice" || contribs[0].Count != 2 {
		t.Errorf("expected Alice with 2 commits, got %s with %d", contribs[0].Name, contribs[0].Count)
	}
}

func TestCollectEmpty(t *testing.T) {
	contribs := Collect(context.Background(), nil, "", "", nil)
	if len(contribs) != 0 {
		t.Errorf("expected 0 contributors, got %d", len(contribs))
	}
}

func TestCollectSingleAuthor(t *testing.T) {
	commits := []*parser.ParsedCommit{
		{Author: "Alice", AuthorEmail: "alice@test.com"},
		{Author: "Alice", AuthorEmail: "alice@test.com"},
		{Author: "Alice", AuthorEmail: "alice@test.com"},
	}

	contribs := Collect(context.Background(), nil, "", "", commits)

	if len(contribs) != 1 {
		t.Errorf("expected 1 contributor, got %d", len(contribs))
	}
	if contribs[0].Count != 3 {
		t.Errorf("expected 3 commits, got %d", contribs[0].Count)
	}
}

func TestCollectSortOrder(t *testing.T) {
	commits := []*parser.ParsedCommit{
		{Author: "Bob", AuthorEmail: "bob@test.com"},
		{Author: "Alice", AuthorEmail: "alice@test.com"},
		{Author: "Bob", AuthorEmail: "bob@test.com"},
		{Author: "Bob", AuthorEmail: "bob@test.com"},
	}

	contribs := Collect(context.Background(), nil, "", "", commits)

	if contribs[0].Name != "Bob" {
		t.Errorf("expected Bob first (3 commits), got %s (%d commits)", contribs[0].Name, contribs[0].Count)
	}
}
