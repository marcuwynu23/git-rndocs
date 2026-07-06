package cmd

import (
	"context"

	"github.com/marcuwynu23/git-rndocs/internal/app"
	"github.com/marcuwynu23/git-rndocs/internal/config"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview release notes in the terminal",
	Long:  `Shows generated release notes in the terminal without writing any files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := &config.GenerateOptions{
			RepoPath: cmd.Flag("repo").Value.String(),
			Version:  cmd.Flag("version").Value.String(),
			From:     cmd.Flag("from").Value.String(),
			To:       cmd.Flag("to").Value.String(),
			Latest:   mustBool(cmd, "latest"),
			Template: cmd.Flag("template").Value.String(),
			Config:   cmd.Flag("config").Value.String(),
			Verbose:  mustBool(cmd, "verbose"),
		}

		application := &app.App{}
		repoPath := opts.RepoPath
		if repoPath == "" {
			repoPath = "."
		}

		if err := application.Init(repoPath, opts); err != nil {
			return err
		}

		return application.Preview(context.Background(), opts)
	},
}

func init() {
	previewCmd.Flags().String("repo", "", "Path to the Git repository")
	previewCmd.Flags().String("version", "", "Specific version/tag to preview")
	previewCmd.Flags().String("from", "", "Starting commit or tag")
	previewCmd.Flags().String("to", "", "Ending commit or tag")
	previewCmd.Flags().Bool("latest", false, "Preview the latest version only")
	previewCmd.Flags().String("template", "", "Template name or path")
	previewCmd.Flags().String("config", "", "Configuration file path")
	previewCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
