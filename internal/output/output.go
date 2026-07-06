package output

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcuwynu23/git-rndocs/internal/github"
	"github.com/marcuwynu23/git-rndocs/internal/releasenotes"
)

type Writer struct {
	outputDir     string
	overwrite     bool
	dryRun        bool
	githubEnabled bool
	githubOpts    *github.ReleaseOptions
}

func NewWriter(outputDir string, overwrite, dryRun bool) *Writer {
	return &Writer{
		outputDir: outputDir,
		overwrite: overwrite,
		dryRun:    dryRun,
	}
}

func (w *Writer) SetGitHubOptions(opts *github.ReleaseOptions) {
	w.githubEnabled = true
	w.githubOpts = opts
}

func (w *Writer) WriteRelease(release *releasenotes.Release) error {
	dir := filepath.Join(w.outputDir, release.Version)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", dir, err)
	}

	filePath := filepath.Join(dir, "RELEASE-NOTES.md")
	if !w.overwrite {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("file already exists: %s (use --overwrite to replace)", filePath)
		}
	}

	if w.dryRun {
		fmt.Printf("[dry-run] Would write: %s\n", filePath)
		fmt.Printf("[dry-run] Content preview:\n%s\n", truncate(release.Markdown, 500))
		return nil
	}

	if err := os.WriteFile(filePath, []byte(release.Markdown), 0644); err != nil {
		return fmt.Errorf("failed to write release notes: %w", err)
	}

	fmt.Printf("Written: %s\n", filePath)

	if w.githubEnabled && w.githubOpts != nil {
		w.githubOpts.Tag = release.TagName
		w.githubOpts.Name = release.Version
		w.githubOpts.Body = release.Markdown

		if w.dryRun {
			fmt.Printf("[dry-run] Would create GitHub release: %s\n", release.Version)
			return nil
		}

		ghRelease, err := github.CreateRelease(nil, w.githubOpts)
		if err != nil {
			return fmt.Errorf("failed to create GitHub release: %w", err)
		}
		fmt.Printf("GitHub Release: %s\n", ghRelease.HTMLURL)
	}

	return nil
}

func (w *Writer) WriteReleases(result *releasenotes.GenerateResult) error {
	for _, release := range result.Releases {
		if err := w.WriteRelease(release); err != nil {
			return fmt.Errorf("failed to write release %s: %w", release.Version, err)
		}
	}
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
