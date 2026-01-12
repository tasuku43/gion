package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	WorkspaceID   string
	WorkspacePath string
	Description   string
}

func List(rootDir string) ([]Entry, []error, error) {
	wsRoot := filepath.Join(rootDir, "workspaces")
	info, err := os.Stat(wsRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("workspaces path is not a directory: %s", wsRoot)
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

		description := ""
		meta, err := LoadMetadata(wsPath)
		if err != nil {
			warnings = append(warnings, fmt.Errorf("workspace %s metadata: %w", wsID, err))
		} else if strings.TrimSpace(meta.Description) != "" {
			description = strings.TrimSpace(meta.Description)
		}

		result := Entry{
			WorkspaceID:   wsID,
			WorkspacePath: wsPath,
			Description:   description,
		}
		results = append(results, result)
	}

	return results, warnings, nil
}
