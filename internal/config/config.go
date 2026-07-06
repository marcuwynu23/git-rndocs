package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Output       string       `mapstructure:"output"`
	Template     string       `mapstructure:"template"`
	ReleaseName  string       `mapstructure:"release_name"`
	Include      []string     `mapstructure:"include"`
	Exclude      []string     `mapstructure:"exclude"`
	GroupUnknown bool         `mapstructure:"group_unknown"`
	GitHub       GitHubConfig `mapstructure:"github"`
	Contributors bool         `mapstructure:"contributors"`
	Statistics   bool         `mapstructure:"statistics"`
}

type GitHubConfig struct {
	Upload bool `mapstructure:"upload"`
}

type GenerateOptions struct {
	RepoPath  string
	OutputDir string
	Version   string
	From      string
	To        string
	All       bool
	Latest    bool
	Template  string
	Config    string
	JSON      bool
	Overwrite bool
	DryRun    bool
	Verbose   bool
}

func DefaultConfig() *Config {
	return &Config{
		Output:       "docs/releases",
		Template:     "default",
		ReleaseName:  "Version {{ .Version }}",
		Include:      []string{"feat", "fix", "perf", "docs", "refactor", "style", "build", "ci", "test", "chore", "revert"},
		Exclude:      []string{},
		GroupUnknown: true,
		GitHub: GitHubConfig{
			Upload: false,
		},
		Contributors: true,
		Statistics:   true,
	}
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(".git-rndocs")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if path != "" {
		v.SetConfigFile(path)
	}

	v.SetDefault("output", "docs/releases")
	v.SetDefault("template", "default")
	v.SetDefault("release_name", "Version {{ .Version }}")
	v.SetDefault("include", []string{"feat", "fix", "perf", "docs", "refactor", "style", "build", "ci", "test", "chore", "revert"})
	v.SetDefault("exclude", []string{})
	v.SetDefault("group_unknown", true)
	v.SetDefault("github.upload", false)
	v.SetDefault("contributors", true)
	v.SetDefault("statistics", true)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func InitConfig(dir string) error {
	cfg := DefaultConfig()
	v := viper.New()
	v.SetConfigType("yaml")

	v.SetDefault("output", cfg.Output)
	v.SetDefault("template", cfg.Template)
	v.SetDefault("release_name", cfg.ReleaseName)
	v.SetDefault("include", cfg.Include)
	v.SetDefault("exclude", cfg.Exclude)
	v.SetDefault("group_unknown", cfg.GroupUnknown)
	v.SetDefault("github.upload", cfg.GitHub.Upload)
	v.SetDefault("contributors", cfg.Contributors)
	v.SetDefault("statistics", cfg.Statistics)

	configPath := filepath.Join(dir, ".git-rndocs.yaml")
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

func CreateDefaultDirectories(dir string) error {
	dirs := []string{
		filepath.Join(dir, "docs", "releases"),
		filepath.Join(dir, "templates"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}
	return nil
}
