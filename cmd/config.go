package cmd

import (
	"fmt"

	"github.com/marcuwynu23/git-rndocs/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and manage git-rndocs configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := cmd.Flag("config").Value.String()

		cfg, err := config.LoadConfig(cfgPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cmd.Flag("get").Value.String() != "" {
			key := cmd.Flag("get").Value.String()
			switch key {
			case "output":
				fmt.Println(cfg.Output)
			case "template":
				fmt.Println(cfg.Template)
			case "release_name":
				fmt.Println(cfg.ReleaseName)
			default:
				return fmt.Errorf("unknown config key: %s", key)
			}
			return nil
		}

		fmt.Printf("output: %s\n", cfg.Output)
		fmt.Printf("template: %s\n", cfg.Template)
		fmt.Printf("release_name: %s\n", cfg.ReleaseName)
		fmt.Printf("contributors: %t\n", cfg.Contributors)
		fmt.Printf("statistics: %t\n", cfg.Statistics)
		fmt.Printf("github.upload: %t\n", cfg.GitHub.Upload)

		return nil
	},
}

func init() {
	configCmd.Flags().String("config", "", "Configuration file path")
	configCmd.Flags().String("get", "", "Get a specific configuration value")
}
