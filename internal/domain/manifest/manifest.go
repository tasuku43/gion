package manifest

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const FileName = "manifest.yaml"

type File struct {
	Version    int                  `yaml:"version"`
	Workspaces map[string]Workspace `yaml:"workspaces"`
}

type Workspace struct {
	Description  string `yaml:"description,omitempty"`
	Mode         string `yaml:"mode"`
	TemplateName string `yaml:"template_name,omitempty"`
	SourceURL    string `yaml:"source_url,omitempty"`
	Repos        []Repo `yaml:"repos"`
}

type Repo struct {
	Alias   string `yaml:"alias"`
	RepoKey string `yaml:"repo_key"`
	Branch  string `yaml:"branch"`
}

func Path(rootDir string) string {
	return filepath.Join(rootDir, FileName)
}

func Load(rootDir string) (File, error) {
	path := Path(rootDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, fmt.Errorf("read manifest: %w", err)
	}
	var file File
	if err := yaml.Unmarshal(data, &file); err != nil {
		return File{}, fmt.Errorf("parse manifest: %w", err)
	}
	if file.Version == 0 {
		file.Version = 1
	}
	if file.Workspaces == nil {
		file.Workspaces = map[string]Workspace{}
	}
	return file, nil
}

func Save(rootDir string, file File) error {
	if file.Version == 0 {
		file.Version = 1
	}
	if file.Workspaces == nil {
		file.Workspaces = map[string]Workspace{}
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(file); err != nil {
		_ = enc.Close()
		return fmt.Errorf("marshal manifest: %w", err)
	}
	if err := enc.Close(); err != nil {
		return fmt.Errorf("close manifest encoder: %w", err)
	}
	if err := os.WriteFile(Path(rootDir), buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}
