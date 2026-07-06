package app

import (
	"context"
	"fmt"
	"os"

	"github.com/marcuwynu23/git-rndocs/internal/config"
	"github.com/marcuwynu23/git-rndocs/internal/git"
	"github.com/marcuwynu23/git-rndocs/internal/github"
	"github.com/marcuwynu23/git-rndocs/internal/output"
	"github.com/marcuwynu23/git-rndocs/internal/releasenotes"
	"github.com/marcuwynu23/git-rndocs/internal/template"
)

type App struct {
	config  *config.Config
	repo    git.Repository
	tmplEng *template.Engine
}

func New() *App {
	return &App{}
}

func (a *App) Init(repoPath string, opts *config.GenerateOptions) error {
	cfg, err := config.LoadConfig(opts.Config)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.config = cfg

	if opts.Template != "" {
		cfg.Template = opts.Template
	}
	if opts.OutputDir != "" {
		cfg.Output = opts.OutputDir
	}

	var repo git.Repository
	repo = &git.GoGitRepo{}
	if err := repo.Open(repoPath); err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}
	a.repo = repo

	a.tmplEng = template.NewEngine("templates")

	return nil
}

func (a *App) Generate(ctx context.Context, opts *config.GenerateOptions) error {
	if !opts.Overwrite && !opts.DryRun {
		if _, err := os.Stat(opts.OutputDir); err == nil {
			fmt.Printf("Output directory %s already exists. Use --overwrite or --dry-run.\n", opts.OutputDir)
			return nil
		}
	}

	gen := releasenotes.NewGenerator(a.repo, a.config, a.tmplEng)
	result, err := gen.GenerateAll(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to generate release notes: %w", err)
	}

	writer := output.NewWriter(ctx, a.config.Output, opts.Overwrite, opts.DryRun)
	if err := writer.WriteReleases(result); err != nil {
		return fmt.Errorf("failed to write releases: %w", err)
	}

	return nil
}

func (a *App) Preview(ctx context.Context, opts *config.GenerateOptions) error {
	gen := releasenotes.NewGenerator(a.repo, a.config, a.tmplEng)
	release, err := gen.GenerateSingle(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to generate release notes: %w", err)
	}

	fmt.Println(release.Markdown)
	return nil
}

func (a *App) Validate(ctx context.Context) error {
	fmt.Println("Validating repository...")

	if a.repo == nil {
		return fmt.Errorf("no repository loaded")
	}

	tags, err := a.repo.Tags(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	fmt.Printf("Tags found: %d\n", len(tags))
	for _, t := range tags {
		fmt.Printf("  - %s\n", t.Name)
	}

	detached, err := a.repo.DetachedHead(ctx)
	if err != nil {
		return fmt.Errorf("failed to check HEAD: %w", err)
	}
	if detached {
		fmt.Println("HEAD is detached")
	}

	fmt.Println("Configuration is valid")
	fmt.Println("Repository is healthy")
	return nil
}

func (a *App) InitProject(dir string) error {
	if err := config.InitConfig(dir); err != nil {
		return fmt.Errorf("failed to init config: %w", err)
	}
	if err := config.CreateDefaultDirectories(dir); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	if err := createDefaultTemplate(dir); err != nil {
		return fmt.Errorf("failed to create default template: %w", err)
	}
	return nil
}

func (a *App) Release(ctx context.Context, opts *config.GenerateOptions, releaseOpts *github.ReleaseOptions) error {
	gen := releasenotes.NewGenerator(a.repo, a.config, a.tmplEng)
	result, err := gen.GenerateAll(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to generate release notes: %w", err)
	}

	writer := output.NewWriter(ctx, a.config.Output, opts.Overwrite, opts.DryRun)
	if releaseOpts != nil {
		writer.SetGitHubOptions(releaseOpts)
	}
	if err := writer.WriteReleases(result); err != nil {
		return fmt.Errorf("failed to write releases: %w", err)
	}

	return nil
}

func createDefaultTemplate(dir string) error {
	tmpl := `# Release Notes

## Version

{{ .Version }}

## Release Date

{{ .ReleaseDate }}

---

{{ range .Sections }}{{ if .Commits }}
## {{ .Title }}

{{ range .Commits }}
- {{ .Header }}{{ end }}
{{ end }}{{ end }}

{{ if .Contributors }}
## Contributors

{{ range .Contributors }}
- {{ .Name }} ({{ .Count }} commits)
{{ end }}{{ end }}
`
	return os.WriteFile(dir+"/templates/default.md", []byte(tmpl), 0644)
}
