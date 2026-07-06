package cmd

import (
	"github.com/marcuwynu23/git-rndocs/internal/app"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize git-rndocs in a project",
	Long: `Creates the default configuration file and directory structure.

Sets up:
  - .git-rndocs.yaml configuration file
  - docs/releases/ directory
  - templates/default.md template`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := cmd.Flag("dir").Value.String()
		if dir == "" {
			dir = "."
		}

		application := &app.App{}
		return application.InitProject(dir)
	},
}

func init() {
	initCmd.Flags().String("dir", "", "Directory to initialize (default: current)")
}
