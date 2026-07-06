package cmd

import (
	"context"

	"github.com/marcuwynu23/git-rndocs/internal/app"
	"github.com/marcuwynu23/git-rndocs/internal/config"
	"github.com/marcuwynu23/git-rndocs/internal/github"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Generate release notes and create a GitHub release",
	Long: `Generates release notes from Git history and optionally uploads them
as a GitHub Release using the GitHub CLI or REST API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := &config.GenerateOptions{
			RepoPath:  cmd.Flag("repo").Value.String(),
			OutputDir: cmd.Flag("output").Value.String(),
			Version:   cmd.Flag("version").Value.String(),
			From:      cmd.Flag("from").Value.String(),
			To:        cmd.Flag("to").Value.String(),
			Latest:    mustBool(cmd, "latest"),
			Template:  cmd.Flag("template").Value.String(),
			Config:    cmd.Flag("config").Value.String(),
			Overwrite: mustBool(cmd, "overwrite"),
			DryRun:    mustBool(cmd, "dry-run"),
			Verbose:   mustBool(cmd, "verbose"),
		}

		upload := mustBool(cmd, "upload")
		draft := mustBool(cmd, "draft")
		prerelease := mustBool(cmd, "prerelease")

		var releaseOpts *github.ReleaseOptions
		if upload {
			releaseOpts = &github.ReleaseOptions{
				Repo:       cmd.Flag("repo").Value.String(),
				Draft:      draft,
				Prerelease: prerelease,
			}
		}

		application := &app.App{}
		repoPath := opts.RepoPath
		if repoPath == "" {
			repoPath = "."
		}

		if err := application.Init(repoPath, opts); err != nil {
			return err
		}

		return application.Release(context.Background(), opts, releaseOpts)
	},
}

func init() {
	releaseCmd.Flags().String("repo", "", "Path to the Git repository")
	releaseCmd.Flags().String("output", "", "Output directory for release notes")
	releaseCmd.Flags().String("version", "", "Specific version/tag to release")
	releaseCmd.Flags().String("from", "", "Starting commit or tag")
	releaseCmd.Flags().String("to", "", "Ending commit or tag")
	releaseCmd.Flags().Bool("latest", true, "Release the latest version")
	releaseCmd.Flags().String("template", "", "Template name or path")
	releaseCmd.Flags().String("config", "", "Configuration file path")
	releaseCmd.Flags().Bool("overwrite", false, "Overwrite existing release notes")
	releaseCmd.Flags().Bool("dry-run", false, "Show what would be generated without writing")
	releaseCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	releaseCmd.Flags().Bool("upload", false, "Upload to GitHub Releases")
	releaseCmd.Flags().Bool("draft", false, "Create a draft release")
	releaseCmd.Flags().Bool("prerelease", false, "Mark as prerelease")
}
