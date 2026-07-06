package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Output != "docs/releases" {
		t.Errorf("expected docs/releases, got %s", cfg.Output)
	}
	if cfg.Template != "default" {
		t.Errorf("expected default template, got %s", cfg.Template)
	}
	if !cfg.Contributors {
		t.Error("expected contributors to be enabled")
	}
	if !cfg.Statistics {
		t.Error("expected statistics to be enabled")
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("expected no error for missing config, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoadConfigExplicitPathNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent config file")
	}
}

func TestInitConfig(t *testing.T) {
	dir := t.TempDir()
	if err := InitConfig(dir); err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	configPath := filepath.Join(dir, ".git-rndocs.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}
}

func TestCreateDefaultDirectories(t *testing.T) {
	dir := t.TempDir()
	if err := CreateDefaultDirectories(dir); err != nil {
		t.Fatalf("CreateDefaultDirectories failed: %v", err)
	}

	releasesDir := filepath.Join(dir, "docs", "releases")
	if _, err := os.Stat(releasesDir); os.IsNotExist(err) {
		t.Fatal("docs/releases directory was not created")
	}
}
