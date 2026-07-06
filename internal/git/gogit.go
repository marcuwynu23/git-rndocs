package git

import (
	"context"
	"fmt"
	"sort"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type GoGitRepo struct {
	repo *gogit.Repository
	path string
}

func (r *GoGitRepo) Open(path string) error {
	repo, err := gogit.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository at %s: %w", path, err)
	}
	r.repo = repo
	r.path = path
	return nil
}

func (r *GoGitRepo) Root() string {
	return r.path
}

func (r *GoGitRepo) Tags(ctx context.Context) ([]Tag, error) {
	tagRefs, err := r.repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	var tags []Tag
	if err := tagRefs.ForEach(func(ref *plumbing.Reference) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		tag := Tag{
			Name: ref.Name().Short(),
		}

		obj, err := r.repo.Object(plumbing.AnyObject, ref.Hash())
		if err != nil {
			return nil
		}

		switch o := obj.(type) {
		case *object.Tag:
			tag.IsAnnotated = true
			tag.Commit = o.Target.String()
			tag.Date = o.Tagger.When
		case *object.Commit:
			tag.Commit = ref.Hash().String()
			tag.Date = o.Committer.When
		}

		tags = append(tags, tag)
		return nil
	}); err != nil {
		return nil, err
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.Before(tags[j].Date)
	})

	return tags, nil
}

func (r *GoGitRepo) CommitsBetween(ctx context.Context, from, to string) ([]Commit, error) {
	toHash, err := r.resolveRevision(to)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve %s: %w", to, err)
	}

	toCommit, err := r.repo.CommitObject(toHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit %s: %w", to, err)
	}

	iter := object.NewCommitPreorderIter(toCommit, nil, nil)
	var commits []Commit

	var fromHash plumbing.Hash
	var hasFrom bool
	if from != "" {
		fromHash, err = r.resolveRevision(from)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s: %w", from, err)
		}
		hasFrom = true
	}

	if err := iter.ForEach(func(c *object.Commit) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if hasFrom && c.Hash == fromHash {
			return storer.ErrStop
		}

		commits = append(commits, Commit{
			Hash:        c.Hash.String(),
			Author:      c.Author.Name,
			AuthorEmail: c.Author.Email,
			Message:     c.Message,
			Date:        c.Author.When,
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return commits, nil
}

func (r *GoGitRepo) CommitCount(ctx context.Context, from, to string) (int, error) {
	commits, err := r.CommitsBetween(ctx, from, to)
	if err != nil {
		return 0, err
	}
	return len(commits), nil
}

func (r *GoGitRepo) DiffStatsBetween(ctx context.Context, from, to string) (*DiffStats, error) {
	fromHash, err := r.resolveRevision(from)
	if err != nil {
		return nil, err
	}
	toHash, err := r.resolveRevision(to)
	if err != nil {
		return nil, err
	}

	fromCommit, err := r.repo.CommitObject(fromHash)
	if err != nil {
		return nil, err
	}
	toCommit, err := r.repo.CommitObject(toHash)
	if err != nil {
		return nil, err
	}

	fromTree, err := fromCommit.Tree()
	if err != nil {
		return nil, err
	}
	toTree, err := toCommit.Tree()
	if err != nil {
		return nil, err
	}

	patch, err := fromTree.Patch(toTree)
	if err != nil {
		return nil, err
	}

	fileStats := patch.Stats()
	stats := &DiffStats{
		FilesChanged: len(fileStats),
	}

	for _, fs := range fileStats {
		stats.Insertions += fs.Addition
		stats.Deletions += fs.Deletion
	}

	return stats, nil
}

func (r *GoGitRepo) Contributors(ctx context.Context, from, to string) ([]string, error) {
	commits, err := r.CommitsBetween(ctx, from, to)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	var contributors []string
	for _, c := range commits {
		name := c.Author
		if name == "" {
			name = c.AuthorEmail
		}
		if _, ok := seen[name]; !ok {
			seen[name] = struct{}{}
			contributors = append(contributors, name)
		}
	}

	sort.Strings(contributors)
	return contributors, nil
}

func (r *GoGitRepo) DetachedHead(ctx context.Context) (bool, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return false, fmt.Errorf("failed to get HEAD: %w", err)
	}
	return ref.Name() == plumbing.HEAD, nil
}

func (r *GoGitRepo) resolveRevision(rev string) (plumbing.Hash, error) {
	if rev == "HEAD" {
		ref, err := r.repo.Head()
		if err != nil {
			return plumbing.ZeroHash, err
		}
		return ref.Hash(), nil
	}

	hash := plumbing.NewHash(rev)
	if hash != plumbing.ZeroHash {
		return hash, nil
	}

	ref, err := r.repo.Reference(plumbing.NewTagReferenceName(rev), true)
	if err != nil {
		ref, err = r.repo.Reference(plumbing.NewBranchReferenceName(rev), true)
		if err != nil {
			return plumbing.ZeroHash, fmt.Errorf("could not resolve %s", rev)
		}
	}
	return ref.Hash(), nil
}

var _ Repository = (*GoGitRepo)(nil)
