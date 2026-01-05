package initcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Result struct {
	RootDir      string
	CreatedDirs  []string
	CreatedFiles []string
	SkippedFiles []string
	SkippedDirs  []string
}

func Run(rootDir string) (Result, error) {
	if rootDir == "" {
		return Result{}, fmt.Errorf("root directory is required")
	}

	result := Result{RootDir: rootDir}

	dirs := []string{
		filepath.Join(rootDir, "bare"),
		filepath.Join(rootDir, "src"),
		filepath.Join(rootDir, "ws"),
	}
	for _, dir := range dirs {
		if exists, err := dirExists(dir); err != nil {
			return Result{}, err
		} else if exists {
			result.SkippedDirs = append(result.SkippedDirs, dir)
			continue
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return Result{}, fmt.Errorf("create dir: %w", err)
		}
		result.CreatedDirs = append(result.CreatedDirs, dir)
	}

	settingsPath := filepath.Join(rootDir, "settings.yaml")
	if exists, err := fileExists(settingsPath); err != nil {
		return Result{}, err
	} else if exists {
		result.SkippedFiles = append(result.SkippedFiles, settingsPath)
	} else {
		if err := writeSettings(settingsPath); err != nil {
			return Result{}, err
		}
		result.CreatedFiles = append(result.CreatedFiles, settingsPath)
	}

	templatesPath := filepath.Join(rootDir, "templates.yaml")
	if exists, err := fileExists(templatesPath); err != nil {
		return Result{}, err
	} else if exists {
		result.SkippedFiles = append(result.SkippedFiles, templatesPath)
	} else {
		if err := writeTemplates(templatesPath); err != nil {
			return Result{}, err
		}
		result.CreatedFiles = append(result.CreatedFiles, templatesPath)
	}

	return result, nil
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

type settingsFile struct {
	Version  int `yaml:"version"`
	Defaults struct {
		BaseRef string `yaml:"base_ref"`
		TTLDays int    `yaml:"ttl_days"`
	} `yaml:"defaults"`
	Repo struct {
		DefaultHost     string `yaml:"default_host"`
		DefaultProtocol string `yaml:"default_protocol"`
	} `yaml:"repo"`
}

func writeSettings(path string) error {
	var file settingsFile
	file.Version = 1
	file.Defaults.BaseRef = ""
	file.Defaults.TTLDays = 30
	file.Repo.DefaultHost = "github.com"
	file.Repo.DefaultProtocol = "https"

	data, err := yaml.Marshal(file)
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write settings: %w", err)
	}
	return nil
}

type templatesFile struct {
	Templates map[string]struct {
		Repos []string `yaml:"repos"`
	} `yaml:"templates"`
}

func writeTemplates(path string) error {
	file := templatesFile{
		Templates: map[string]struct {
			Repos []string `yaml:"repos"`
		}{},
	}
	data, err := yaml.Marshal(file)
	if err != nil {
		return fmt.Errorf("marshal templates: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write templates: %w", err)
	}
	return nil
}
