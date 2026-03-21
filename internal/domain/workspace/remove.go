package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tasuku43/gion/internal/infra/gitcmd"
	"github.com/tasuku43/gion/internal/infra/paths"
	"github.com/tasuku43/gion/internal/infra/shellaction"
)

type RemoveOptions struct {
	AllowStatusError bool
	AllowDirty       bool
}

func Remove(ctx context.Context, rootDir, workspaceID string) error {
	return RemoveWithOptions(ctx, rootDir, workspaceID, RemoveOptions{})
}

func RemoveWithOptions(ctx context.Context, rootDir, workspaceID string, opts RemoveOptions) error {
	if workspaceID == "" {
		return fmt.Errorf("workspace id is required")
	}
	if rootDir == "" {
		return fmt.Errorf("root directory is required")
	}
	if err := validateWorkspaceID(ctx, workspaceID); err != nil {
		return err
	}

	wsDir := WorkspaceDir(rootDir, workspaceID)
	if exists, err := paths.DirExists(wsDir); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("workspace does not exist: %s", wsDir)
	}
	shifted, err := shiftCWDToRootIfInsideWorkspace(rootDir, wsDir)
	if err != nil {
		return err
	}

	repos, warnings, err := ScanRepos(ctx, wsDir)
	if err != nil {
		return err
	}
	_ = warnings

	for _, repo := range repos {
		if repo.WorktreePath == "" {
			return fmt.Errorf("missing worktree path for alias %q", repo.Alias)
		}
		statusOut, statusErr := gitStatusPorcelain(ctx, repo.WorktreePath)
		if statusErr != nil {
			if !opts.AllowStatusError {
				return fmt.Errorf("check status for %q: %w", repo.Alias, statusErr)
			}
			continue
		}
		_, _, _, _, _, dirty, _, _, _, _, _, _ := parseStatusPorcelainV2(statusOut, "")
		if dirty {
			if !opts.AllowDirty {
				return fmt.Errorf("workspace has dirty changes: %s", repo.Alias)
			}
		}
	}

	for _, repo := range repos {
		if repo.StorePath == "" {
			continue
		}
		if repo.WorktreePath == "" {
			return fmt.Errorf("missing worktree path for alias %q", repo.Alias)
		}
		force := opts.AllowDirty
		if force {
			gitcmd.Logf("git worktree remove --force %s", repo.WorktreePath)
		} else {
			gitcmd.Logf("git worktree remove %s", repo.WorktreePath)
		}
		if err := gitcmd.WorktreeRemove(ctx, repo.StorePath, repo.WorktreePath, force); err != nil {
			return fmt.Errorf("remove worktree %q: %w", repo.Alias, err)
		}
	}

	if err := os.RemoveAll(wsDir); err != nil {
		return fmt.Errorf("remove workspace dir: %w", err)
	}
	if shifted {
		if err := shellaction.EmitCD(rootDir); err != nil {
			return err
		}
	}

	return nil
}

func shiftCWDToRootIfInsideWorkspace(rootDir, workspaceDir string) (bool, error) {
	if strings.TrimSpace(rootDir) == "" || strings.TrimSpace(workspaceDir) == "" {
		return false, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("get working dir: %w", err)
	}
	if !pathInside(workspaceDir, cwd) {
		return false, nil
	}
	if err := os.Chdir(rootDir); err != nil {
		return false, fmt.Errorf("shift process cwd to root: %w", err)
	}
	return true, nil
}

func pathInside(base, target string) bool {
	if strings.TrimSpace(base) == "" || strings.TrimSpace(target) == "" {
		return false
	}
	basePath := filepath.Clean(base)
	targetPath := filepath.Clean(target)
	if resolved, err := filepath.EvalSymlinks(basePath); err == nil {
		basePath = filepath.Clean(resolved)
	}
	if resolved, err := filepath.EvalSymlinks(targetPath); err == nil {
		targetPath = filepath.Clean(resolved)
	}
	rel, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return false
	}
	rel = filepath.Clean(rel)
	if rel == "." {
		return true
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
