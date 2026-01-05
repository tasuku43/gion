package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingReturnsDefault(t *testing.T) {
	temp := t.TempDir()
	t.Setenv("HOME", temp)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Defaults.BaseRef != "" {
		t.Fatalf("expected default base_ref to be empty, got %q", cfg.Defaults.BaseRef)
	}
	if cfg.Paths.WsDir != "ws" {
		t.Fatalf("expected ws dir default, got %q", cfg.Paths.WsDir)
	}
}

func TestLoadConfigRoot(t *testing.T) {
	temp := t.TempDir()
	t.Setenv("HOME", temp)
	configDir := filepath.Join(temp, ".config", "gws")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	data := []byte("root: /tmp/custom-root\n")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Root != "/tmp/custom-root" {
		t.Fatalf("expected /tmp/custom-root, got %s", cfg.Root)
	}
	if cfg.Paths.ReposDir != "bare" {
		t.Fatalf("expected repos dir default, got %q", cfg.Paths.ReposDir)
	}
}
