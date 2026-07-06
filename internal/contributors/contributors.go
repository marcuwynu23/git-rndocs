package contributors

import (
	"context"
	"sort"

	"github.com/marcuwynu23/git-rndocs/internal/git"
	"github.com/marcuwynu23/git-rndocs/internal/parser"
)

type Contributor struct {
	Name   string `json:"name"`
	Count  int    `json:"commit_count"`
}

func Collect(ctx context.Context, repo git.Repository, from, to string, commits []*parser.ParsedCommit) []Contributor {
	counts := make(map[string]int)

	for _, c := range commits {
		name := c.Author
		if name == "" {
			name = c.AuthorEmail
		}
		if name == "" {
			name = "Unknown"
		}
		counts[name]++
	}

	var contributors []Contributor
	for name, count := range counts {
		contributors = append(contributors, Contributor{Name: name, Count: count})
	}

	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Count > contributors[j].Count
	})

	return contributors
}
