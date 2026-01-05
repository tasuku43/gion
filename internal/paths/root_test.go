package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveRootFlagOverrides(t *testing.T) {
	t.Setenv("GWS_ROOT", "/tmp/ignore")
	root, err := ResolveRoot("/tmp/custom")
	if err != nil {
		t.Fatalf("ResolveRoot error: %v", err)
	}
	if root != "/tmp/custom" {
		t.Fatalf("expected /tmp/custom, got %s", root)
	}
}

func TestResolveRootEnvOverridesConfig(t *testing.T) {
	t.Setenv("GWS_ROOT", "/tmp/env-root")
	root, err := ResolveRoot("")
	if err != nil {
		t.Fatalf("ResolveRoot error: %v", err)
	}
	if root != "/tmp/env-root" {
		t.Fatalf("expected /tmp/env-root, got %s", root)
	}
}

func TestResolveRootConfig(t *testing.T) {
	temp := t.TempDir()
	t.Setenv("HOME", temp)
	configDir := filepath.Join(temp, ".config", "gws")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("root: /tmp/config-root\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	root, err := ResolveRoot("")
	if err != nil {
		t.Fatalf("ResolveRoot error: %v", err)
	}
	if root != "/tmp/config-root" {
		t.Fatalf("expected /tmp/config-root, got %s", root)
	}
}

func TestResolveRootDefault(t *testing.T) {
	temp := t.TempDir()
	t.Setenv("HOME", temp)
	root, err := ResolveRoot("")
	if err != nil {
		t.Fatalf("ResolveRoot error: %v", err)
	}
	expected := filepath.Join(temp, "gws")
	if root != expected {
		t.Fatalf("expected %s, got %s", expected, root)
	}
}
