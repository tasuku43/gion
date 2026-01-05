package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  int            `yaml:"version"`
	Root     string         `yaml:"root"`
	Paths    PathsConfig    `yaml:"paths"`
	Defaults DefaultsConfig `yaml:"defaults"`
	Naming   NamingConfig   `yaml:"naming"`
	Repo     RepoConfig     `yaml:"repo"`
}

type PathsConfig struct {
	ReposDir string `yaml:"repos_dir"`
	WsDir    string `yaml:"ws_dir"`
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

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gws", "config.yaml"), nil
}

func DefaultConfig() Config {
	return Config{
		Version: 1,
		Root:    "",
		Paths: PathsConfig{
			ReposDir: "bare",
			WsDir:    "ws",
		},
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

func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		var err error
		path, err = DefaultPath()
		if err != nil {
			return Config{}, err
		}
	}

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
