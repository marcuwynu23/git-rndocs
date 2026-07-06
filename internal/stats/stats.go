package stats

import (
	"context"
	"fmt"

	"github.com/marcuwynu23/git-rndocs/internal/git"
	"github.com/marcuwynu23/git-rndocs/internal/parser"
)

type Statistics struct {
	TotalCommits      int                    `json:"total_commits"`
	TotalContributors int                    `json:"total_contributors"`
	FilesChanged      int                    `json:"files_changed"`
	Insertions        int                    `json:"insertions"`
	Deletions         int                    `json:"deletions"`
	ByCategory        map[string]int         `json:"by_category"`
}

func Collect(ctx context.Context, repo git.Repository, from, to string, categorized map[parser.CommitType][]*parser.ParsedCommit) (*Statistics, error) {
	stats := &Statistics{
		ByCategory: make(map[string]int),
	}

	for ct, commits := range categorized {
		stats.ByCategory[ct.DisplayName()] = len(commits)
		stats.TotalCommits += len(commits)
	}

	if repo != nil {
		contributors, err := repo.Contributors(ctx, from, to)
		if err != nil {
			return nil, fmt.Errorf("failed to get contributors: %w", err)
		}
		stats.TotalContributors = len(contributors)

		diffStats, err := repo.DiffStatsBetween(ctx, from, to)
		if err != nil {
			return nil, fmt.Errorf("failed to get diff stats: %w", err)
		}
		if diffStats != nil {
			stats.FilesChanged = diffStats.FilesChanged
			stats.Insertions = diffStats.Insertions
			stats.Deletions = diffStats.Deletions
		}
	}

	return stats, nil
}
