package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	// Create first commit
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	_, err = worktree.Add("README.md")
	if err != nil {
		t.Fatalf("failed to stage: %v", err)
	}
	_, err = worktree.Commit("feat: initial commit", &gogit.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@test.com"},
	})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Tag v1.0.0
	head, err := repo.Head()
	if err != nil {
		t.Fatalf("failed to get HEAD: %v", err)
	}
	_, err = repo.CreateTag("v1.0.0", head.Hash(), nil)
	if err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	// Create second commit
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	_, err = worktree.Add("main.go")
	if err != nil {
		t.Fatalf("failed to stage: %v", err)
	}
	_, err = worktree.Commit("feat(core): add main package\n\nBREAKING CHANGE: new API", &gogit.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@test.com"},
	})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Tag v2.0.0
	head, err = repo.Head()
	if err != nil {
		t.Fatalf("failed to get HEAD: %v", err)
	}
	_, err = repo.CreateTag("v2.0.0", head.Hash(), nil)
	if err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	return dir
}

func TestOpenAndTags(t *testing.T) {
	dir := setupTestRepo(t)

	var r GoGitRepo
	if err := r.Open(dir); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	tags, err := r.Tags(context.Background())
	if err != nil {
		t.Fatalf("Tags failed: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}

	if tags[0].Name != "v1.0.0" {
		t.Errorf("expected first tag v1.0.0, got %s", tags[0].Name)
	}
	if tags[1].Name != "v2.0.0" {
		t.Errorf("expected second tag v2.0.0, got %s", tags[1].Name)
	}
}

func TestCommitsBetween(t *testing.T) {
	dir := setupTestRepo(t)

	var r GoGitRepo
	if err := r.Open(dir); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	commits, err := r.CommitsBetween(context.Background(), "v1.0.0", "v2.0.0")
	if err != nil {
		t.Fatalf("CommitsBetween failed: %v", err)
	}

	if len(commits) != 1 {
		t.Fatalf("expected 1 commit between v1.0.0 and v2.0.0, got %d", len(commits))
	}

	if commits[0].Author != "Test" {
		t.Errorf("expected author Test, got %s", commits[0].Author)
	}
}

func TestCommitsFromBeginning(t *testing.T) {
	dir := setupTestRepo(t)

	var r GoGitRepo
	if err := r.Open(dir); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	commits, err := r.CommitsBetween(context.Background(), "", "v1.0.0")
	if err != nil {
		t.Fatalf("CommitsBetween failed: %v", err)
	}

	if len(commits) != 1 {
		t.Fatalf("expected 1 commit from beginning to v1.0.0, got %d", len(commits))
	}
}

func TestDiffStatsBetween(t *testing.T) {
	dir := setupTestRepo(t)

	var r GoGitRepo
	if err := r.Open(dir); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	stats, err := r.DiffStatsBetween(context.Background(), "v1.0.0", "v2.0.0")
	if err != nil {
		t.Fatalf("DiffStatsBetween failed: %v", err)
	}

	if stats.FilesChanged == 0 {
		t.Error("expected at least 1 file changed")
	}
}

func TestContributors(t *testing.T) {
	dir := setupTestRepo(t)

	var r GoGitRepo
	if err := r.Open(dir); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	contribs, err := r.Contributors(context.Background(), "v1.0.0", "HEAD")
	if err != nil {
		t.Fatalf("Contributors failed: %v", err)
	}

	if len(contribs) == 0 {
		t.Error("expected at least 1 contributor")
	}
}

func TestDetachedHead(t *testing.T) {
	dir := setupTestRepo(t)

	var r GoGitRepo
	if err := r.Open(dir); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	detached, err := r.DetachedHead(context.Background())
	if err != nil {
		t.Fatalf("DetachedHead failed: %v", err)
	}

	if detached {
		t.Error("expected HEAD not to be detached")
	}
}

func TestRoot(t *testing.T) {
	dir := setupTestRepo(t)

	var r GoGitRepo
	if err := r.Open(dir); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if r.Root() != dir {
		t.Errorf("expected root %s, got %s", dir, r.Root())
	}
}

func TestCommitCount(t *testing.T) {
	dir := setupTestRepo(t)

	var r GoGitRepo
	if err := r.Open(dir); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	count, err := r.CommitCount(context.Background(), "v1.0.0", "v2.0.0")
	if err != nil {
		t.Fatalf("CommitCount failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 commit, got %d", count)
	}
}
