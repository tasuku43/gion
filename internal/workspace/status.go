package workspace

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tasuku43/gws/internal/gitcmd"
)

type StatusResult struct {
	WorkspaceID string
	Repos       []RepoStatus
}

type RepoStatus struct {
	Alias          string
	Branch         string
	Head           string
	Dirty          bool
	UntrackedCount int
	StagedCount    int
	UnstagedCount  int
	UnmergedCount  int
	WorktreePath   string
	RawStatus      string
	Error          error
}

func Status(ctx context.Context, rootDir, workspaceID string) (StatusResult, error) {
	if workspaceID == "" {
		return StatusResult{}, fmt.Errorf("workspace id is required")
	}
	if rootDir == "" {
		return StatusResult{}, fmt.Errorf("root directory is required")
	}

	wsDir := filepath.Join(rootDir, "ws", workspaceID)
	if exists, err := pathExists(wsDir); err != nil {
		return StatusResult{}, err
	} else if !exists {
		return StatusResult{}, fmt.Errorf("workspace does not exist: %s", wsDir)
	}

	manifestPath := filepath.Join(wsDir, manifestDirName, manifestFileName)
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return StatusResult{}, err
	}

	result := StatusResult{
		WorkspaceID: workspaceID,
	}
	for _, repo := range manifest.Repos {
		repoStatus := RepoStatus{
			Alias:        repo.Alias,
			Branch:       repo.Branch,
			WorktreePath: repo.WorktreePath,
		}

		statusOut, statusErr := gitStatusPorcelain(ctx, repo.WorktreePath)
		if statusErr != nil {
			repoStatus.Error = statusErr
			result.Repos = append(result.Repos, repoStatus)
			continue
		}

		repoStatus.RawStatus = statusOut
		repoStatus.Branch, repoStatus.Head, repoStatus.Dirty, repoStatus.UntrackedCount, repoStatus.StagedCount, repoStatus.UnstagedCount, repoStatus.UnmergedCount = parseStatusPorcelainV2(statusOut, repoStatus.Branch)
		result.Repos = append(result.Repos, repoStatus)
	}

	return result, nil
}

func gitStatusPorcelain(ctx context.Context, worktreePath string) (string, error) {
	res, err := gitcmd.Run(ctx, []string{"status", "--porcelain=v2", "-b"}, gitcmd.Options{Dir: worktreePath})
	if err != nil {
		if strings.TrimSpace(res.Stderr) != "" {
			return "", fmt.Errorf("git status failed: %w: %s", err, strings.TrimSpace(res.Stderr))
		}
		return "", fmt.Errorf("git status failed: %w", err)
	}
	return res.Stdout, nil
}

func parseStatusPorcelainV2(output, fallbackBranch string) (string, string, bool, int, int, int, int) {
	branch := fallbackBranch
	var head string
	var dirty bool
	var untracked int
	var staged int
	var unstaged int
	var unmerged int

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			switch fields[1] {
			case "branch.oid":
				if fields[2] != "(initial)" {
					head = shortSHA(fields[2])
				}
			case "branch.head":
				if fields[2] != "(detached)" && fields[2] != "(unknown)" {
					branch = fields[2]
				}
			}
			continue
		}

		if strings.HasPrefix(line, "? ") {
			untracked++
			dirty = true
			continue
		}

		if strings.HasPrefix(line, "u ") {
			unmerged++
			dirty = true
			continue
		}
		if strings.HasPrefix(line, "1 ") || strings.HasPrefix(line, "2 ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				xy := fields[1]
				if len(xy) >= 2 {
					if xy[0] != '.' {
						staged++
					}
					if xy[1] != '.' {
						unstaged++
					}
					if xy[0] != '.' || xy[1] != '.' {
						dirty = true
					}
				}
			}
			continue
		}
		dirty = true
	}

	return branch, head, dirty, untracked, staged, unstaged, unmerged
}

func shortSHA(oid string) string {
	if len(oid) <= 7 {
		return oid
	}
	return oid[:7]
}
