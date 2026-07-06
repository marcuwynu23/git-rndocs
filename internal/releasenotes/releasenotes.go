package releasenotes

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/marcuwynu23/git-rndocs/internal/config"
	"github.com/marcuwynu23/git-rndocs/internal/contributors"
	"github.com/marcuwynu23/git-rndocs/internal/git"
	"github.com/marcuwynu23/git-rndocs/internal/markdown"
	"github.com/marcuwynu23/git-rndocs/internal/parser"
	"github.com/marcuwynu23/git-rndocs/internal/stats"
	"github.com/marcuwynu23/git-rndocs/internal/template"
)

type Generator struct {
	repo     git.Repository
	cfg      *config.Config
	tmplEng  *template.Engine
}

func NewGenerator(repo git.Repository, cfg *config.Config, tmplEng *template.Engine) *Generator {
	return &Generator{
		repo:    repo,
		cfg:     cfg,
		tmplEng: tmplEng,
	}
}

type Release struct {
	Version     string
	TagName     string
	From        string
	To          string
	Commits     []*parser.ParsedCommit
	Categorized map[parser.CommitType][]*parser.ParsedCommit
	Breaking    []*parser.ParsedCommit
	Stats       *stats.Statistics
	Contributors []contributors.Contributor
	Markdown    string
}

type GenerateResult struct {
	Releases []*Release
}

func (g *Generator) GenerateAll(ctx context.Context, opts *config.GenerateOptions) (*GenerateResult, error) {
	tags, err := g.repo.Tags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	if len(tags) == 0 {
		// No tags, generate from beginning of history
		release, err := g.generateSingle(ctx, "", "HEAD", opts)
		if err != nil {
			return nil, err
		}
		return &GenerateResult{Releases: []*Release{release}}, nil
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.Before(tags[j].Date)
	})

	var releases []*Release

	if opts.Latest {
		latest := tags[len(tags)-1]
		from := latest.Name
		to := "HEAD"
		if !opts.All {
			from = ""
		}
		release, err := g.generateSingle(ctx, from, to, opts)
		if err != nil {
			return nil, err
		}
		return &GenerateResult{Releases: []*Release{release}}, nil
	}

	for i := 0; i < len(tags); i++ {
		from := ""
		if i > 0 {
			from = tags[i-1].Name
		}
		to := tags[i].Name

		release, err := g.generateSingle(ctx, from, to, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate for %s: %w", to, err)
		}
		releases = append(releases, release)
	}

	lastTag := tags[len(tags)-1].Name
	if !opts.All && opts.Version == "" {
		release, err := g.generateSingle(ctx, lastTag, "HEAD", opts)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}

	return &GenerateResult{Releases: releases}, nil
}

func (g *Generator) GenerateSingle(ctx context.Context, opts *config.GenerateOptions) (*Release, error) {
	from := opts.From
	to := opts.To
	if to == "" {
		to = "HEAD"
	}

	if opts.Version != "" {
		from = opts.Version
		return g.generateSingle(ctx, from, to, opts)
	}

	tags, err := g.repo.Tags(ctx)
	if err != nil {
		return nil, err
	}

	if len(tags) > 0 {
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Date.Before(tags[j].Date)
		})
		from = tags[len(tags)-1].Name
	}

	return g.generateSingle(ctx, from, to, opts)
}

func (g *Generator) generateSingle(ctx context.Context, from, to string, opts *config.GenerateOptions) (*Release, error) {
	commits, err := g.repo.CommitsBetween(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	var parsedCommits []*parser.ParsedCommit
	for _, c := range commits {
		pc := parser.ParseCommit(c.Message, c.Hash, c.Author, c.AuthorEmail)
		parsedCommits = append(parsedCommits, pc)
	}

	categorized := parser.CategorizeCommits(parsedCommits, g.cfg.Include, g.cfg.Exclude)

	var breaking []*parser.ParsedCommit
	for _, c := range parsedCommits {
		if c.Breaking {
			breaking = append(breaking, c)
		}
	}

	var version string
	if from == "" {
		version = "initial"
	} else if to == "HEAD" {
		version = fmt.Sprintf("%s-next", from)
	} else {
		version = to
	}

	var stat *stats.Statistics
	if g.cfg.Statistics {
		s, err := stats.Collect(ctx, g.repo, from, to, categorized)
		if err == nil {
			stat = s
		}
	}

	var contribs []contributors.Contributor
	if g.cfg.Contributors {
		contribs = contributors.Collect(ctx, g.repo, from, to, parsedCommits)
	}

	date := time.Now().UTC().Format("2006-01-02")

	tmplName := opts.Template
	if tmplName == "" {
		tmplName = g.cfg.Template
	}

	var md string
	if strings.HasPrefix(tmplName, "markdown:") || tmplName == "default" {
		tmplName = strings.TrimPrefix(tmplName, "markdown:")
		rd := &markdown.ReleaseData{
			Version:      version,
			Date:         date,
			Sections:     buildSections(categorized),
			Breaking:     breaking,
			Contributors: contribs,
			Statistics:   stat,
			From:         from,
			To:           to,
		}
		md = markdown.BuildMarkdown(rd)
	} else {
		rd := g.tmplEng.BuildRenderData(version, date, "", from, to, categorized, breaking, contribs, stat)
		rendered, err := g.tmplEng.Render(tmplName, rd)
		if err != nil {
			return nil, fmt.Errorf("template rendering failed: %w", err)
		}
		md = rendered
	}

	return &Release{
		Version:      version,
		TagName:      version,
		From:         from,
		To:           to,
		Commits:      parsedCommits,
		Categorized:  categorized,
		Breaking:     breaking,
		Stats:        stat,
		Contributors: contribs,
		Markdown:     md,
	}, nil
}

func buildSections(categorized map[parser.CommitType][]*parser.ParsedCommit) []markdown.Section {
	var sections []markdown.Section
	order := []parser.CommitType{
		parser.TypeFeat, parser.TypeFix, parser.TypePerf, parser.TypeDocs,
		parser.TypeRefactor, parser.TypeStyle, parser.TypeBuild, parser.TypeCI,
		parser.TypeTest, parser.TypeChore, parser.TypeRevert, parser.TypeUnknown,
	}
	for _, ct := range order {
		if commits, ok := categorized[ct]; ok && len(commits) > 0 {
			sections = append(sections, markdown.Section{
				Title:   ct.DisplayName(),
				Commits: commits,
			})
		}
	}
	return sections
}
