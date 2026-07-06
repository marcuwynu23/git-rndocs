package git

import (
	"context"
	"time"
)

type Commit struct {
	Hash      string
	Author    string
	AuthorEmail string
	Message   string
	Date      time.Time
}

type Tag struct {
	Name      string
	Commit    string
	IsAnnotated bool
	Date      time.Time
}

type DiffStats struct {
	FilesChanged int
	Insertions   int
	Deletions    int
}

type Repository interface {
	Open(path string) error
	Tags(ctx context.Context) ([]Tag, error)
	CommitsBetween(ctx context.Context, from, to string) ([]Commit, error)
	CommitCount(ctx context.Context, from, to string) (int, error)
	DiffStatsBetween(ctx context.Context, from, to string) (*DiffStats, error)
	Contributors(ctx context.Context, from, to string) ([]string, error)
	DetachedHead(ctx context.Context) (bool, error)
	Root() string
}
