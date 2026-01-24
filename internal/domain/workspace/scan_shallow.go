package workspace

import (
	"context"
	"os"
	"path/filepath"
)

// ScanReposShallow lists repositories under a workspace by checking for .git presence first.
// It still uses git commands to extract metadata once a .git entry is found.
func ScanReposShallow(ctx context.Context, wsDir string) ([]Repo, []error, error) {
	entries, err := os.ReadDir(wsDir)
	if err != nil {
		return nil, nil, err
	}
	var repos []Repo
	var warnings []error
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == MetadataDirName {
			continue
		}
		repoPath := filepath.Join(wsDir, entry.Name())
		if !hasGitEntry(repoPath) {
			continue
		}
		repo, warn, ok := inspectRepo(ctx, repoPath, entry.Name())
		if !ok {
			if warn != nil {
				warnings = append(warnings, warn)
			}
			continue
		}
		if warn != nil {
			warnings = append(warnings, warn)
		}
		repos = append(repos, repo)
	}
	return repos, warnings, nil
}

func hasGitEntry(repoPath string) bool {
	_, err := os.Stat(filepath.Join(repoPath, ".git"))
	return err == nil
}
