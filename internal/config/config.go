package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  int            `yaml:"version"`
	Defaults DefaultsConfig `yaml:"defaults"`
	Naming   NamingConfig   `yaml:"naming"`
	Repo     RepoConfig     `yaml:"repo"`
}

type DefaultsConfig struct {
	BaseRef string `yaml:"base_ref"`
	TTLDays int    `yaml:"ttl_days"`
}

type NamingConfig struct {
	WorkspaceIDMustBeValidRefname bool `yaml:"workspace_id_must_be_valid_refname"`
	BranchEqualsWorkspaceID       bool `yaml:"branch_equals_workspace_id"`
}

type RepoConfig struct {
	DefaultHost     string `yaml:"default_host"`
	DefaultProtocol string `yaml:"default_protocol"`
}

func DefaultConfig() Config {
	return Config{
		Version: 1,
		Defaults: DefaultsConfig{
			BaseRef: "",
			TTLDays: 30,
		},
		Naming: NamingConfig{
			WorkspaceIDMustBeValidRefname: true,
			BranchEqualsWorkspaceID:       true,
		},
		Repo: RepoConfig{
			DefaultHost:     "github.com",
			DefaultProtocol: "https",
		},
	}
}

func Load(rootDir string) (Config, error) {
	cfg := DefaultConfig()

	if rootDir == "" {
		return Config{}, fmt.Errorf("root directory is required")
	}

	path := filepath.Join(rootDir, "settings.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
