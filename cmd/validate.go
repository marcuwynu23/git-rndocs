package cmd

import (
	"context"

	"github.com/marcuwynu23/git-rndocs/internal/app"
	"github.com/marcuwynu23/git-rndocs/internal/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate repository and configuration",
	Long: `Checks the Git repository and configuration for issues.

Validates tags, Git history, configuration, and templates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := cmd.Flag("repo").Value.String()
		if repoPath == "" {
			repoPath = "."
		}

		application := &app.App{}
		if err := application.Init(repoPath, &config.GenerateOptions{}); err != nil {
			return err
		}

		return application.Validate(context.Background())
	},
}

func init() {
	validateCmd.Flags().String("repo", "", "Path to the Git repository")
}
