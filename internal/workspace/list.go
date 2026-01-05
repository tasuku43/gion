package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

type Entry struct {
	WorkspaceID   string
	WorkspacePath string
	ManifestPath  string
	Manifest      *Manifest
	Warning       error
}

func List(rootDir string) ([]Entry, []error, error) {
	wsRoot := filepath.Join(rootDir, "ws")
	info, err := os.Stat(wsRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("ws path is not a directory: %s", wsRoot)
	}

	entries, err := os.ReadDir(wsRoot)
	if err != nil {
		return nil, nil, err
	}

	var results []Entry
	var warnings []error

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		wsID := entry.Name()
		wsPath := filepath.Join(wsRoot, wsID)
		manifestPath := filepath.Join(wsPath, manifestDirName, manifestFileName)

		result := Entry{
			WorkspaceID:   wsID,
			WorkspacePath: wsPath,
			ManifestPath:  manifestPath,
		}

		manifest, err := LoadManifest(manifestPath)
		if err != nil {
			result.Warning = err
			warnings = append(warnings, fmt.Errorf("workspace %s: %w", wsID, err))
			results = append(results, result)
			continue
		}
		result.Manifest = &manifest
		results = append(results, result)
	}

	return results, warnings, nil
}
