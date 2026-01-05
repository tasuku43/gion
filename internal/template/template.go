package template

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const FileName = "templates.yaml"

type File struct {
	Templates map[string]Template `yaml:"templates"`
}

type Template struct {
	Repos []string `yaml:"repos"`
}

func Load(rootDir string) (File, error) {
	if rootDir == "" {
		return File{}, fmt.Errorf("root directory is required")
	}
	path := filepath.Join(rootDir, FileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return File{}, fmt.Errorf("templates file not found: %s", path)
		}
		return File{}, err
	}

	var file File
	if err := yaml.Unmarshal(data, &file); err != nil {
		return File{}, err
	}
	if file.Templates == nil {
		file.Templates = map[string]Template{}
	}
	return file, nil
}

func Save(rootDir string, file File) error {
	if rootDir == "" {
		return fmt.Errorf("root directory is required")
	}
	path := filepath.Join(rootDir, FileName)
	data, err := yaml.Marshal(file)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	return nil
}

func Names(file File) []string {
	var names []string
	for name := range file.Templates {
		if strings.TrimSpace(name) == "" {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
