package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingReturnsDefault(t *testing.T) {
	rootDir := t.TempDir()
	cfg, err := Load(rootDir)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Defaults.BaseRef != "" {
		t.Fatalf("expected default base_ref to be empty, got %q", cfg.Defaults.BaseRef)
	}
	if cfg.Repo.DefaultHost != "github.com" {
		t.Fatalf("expected default host, got %q", cfg.Repo.DefaultHost)
	}
}

func TestLoadConfigRoot(t *testing.T) {
	rootDir := t.TempDir()
	configPath := filepath.Join(rootDir, "settings.yaml")
	data := []byte("defaults:\n  ttl_days: 10\nrepo:\n  default_host: example.com\n")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(rootDir)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Defaults.TTLDays != 10 {
		t.Fatalf("expected ttl_days=10, got %d", cfg.Defaults.TTLDays)
	}
	if cfg.Repo.DefaultHost != "example.com" {
		t.Fatalf("expected default_host=example.com, got %q", cfg.Repo.DefaultHost)
	}
}
