package cmd

import (
	"context"

	"github.com/marcuwynu23/git-rndocs/internal/app"
	"github.com/marcuwynu23/git-rndocs/internal/config"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate release notes",
	Long: `Generate professional release notes from Git history.

Analyzes commits between versions, categorizes them by type, and generates
structured Markdown release notes. Supports Conventional Commits.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := &config.GenerateOptions{
			RepoPath:  cmd.Flag("repo").Value.String(),
			OutputDir: cmd.Flag("output").Value.String(),
			Version:   cmd.Flag("version").Value.String(),
			From:      cmd.Flag("from").Value.String(),
			To:        cmd.Flag("to").Value.String(),
			All:       mustBool(cmd, "all"),
			Latest:    mustBool(cmd, "latest"),
			Template:  cmd.Flag("template").Value.String(),
			Config:    cmd.Flag("config").Value.String(),
			JSON:      mustBool(cmd, "json"),
			Overwrite: mustBool(cmd, "overwrite"),
			DryRun:    mustBool(cmd, "dry-run"),
			Verbose:   mustBool(cmd, "verbose"),
		}

		application := &app.App{}
		repoPath := opts.RepoPath
		if repoPath == "" {
			repoPath = "."
		}

		if err := application.Init(repoPath, opts); err != nil {
			return err
		}

		return application.Generate(context.Background(), opts)
	},
}

func init() {
	generateCmd.Flags().String("repo", "", "Path to the Git repository")
	generateCmd.Flags().String("output", "", "Output directory for release notes")
	generateCmd.Flags().String("version", "", "Specific version/tag to generate notes for")
	generateCmd.Flags().String("from", "", "Starting commit or tag")
	generateCmd.Flags().String("to", "", "Ending commit or tag")
	generateCmd.Flags().Bool("all", false, "Generate release notes for all versions")
	generateCmd.Flags().Bool("latest", false, "Generate release notes for the latest version only")
	generateCmd.Flags().String("template", "", "Template name or path")
	generateCmd.Flags().String("config", "", "Configuration file path")
	generateCmd.Flags().Bool("json", false, "Output in JSON format")
	generateCmd.Flags().Bool("overwrite", false, "Overwrite existing release notes")
	generateCmd.Flags().Bool("dry-run", false, "Show what would be generated without writing")
	generateCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}

func mustBool(cmd *cobra.Command, name string) bool {
	v, _ := cmd.Flags().GetBool(name)
	return v
}
