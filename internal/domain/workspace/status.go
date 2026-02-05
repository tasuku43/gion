package workspace

import (
	"context"
	"fmt"

	coregitparse "github.com/tasuku43/gion-core/gitparse"
	"github.com/tasuku43/gion/internal/infra/gitcmd"
	"github.com/tasuku43/gion/internal/infra/paths"
)

type StatusResult struct {
	WorkspaceID string
	Repos       []RepoStatus
	Warnings    []error
}

type RepoStatus struct {
	Alias          string
	Branch         string
	Upstream       string
	Head           string
	Detached       bool
	HeadMissing    bool
	Dirty          bool
	UntrackedCount int
	StagedCount    int
	UnstagedCount  int
	UnmergedCount  int
	AheadCount     int
	BehindCount    int
	WorktreePath   string
	RawStatus      string
	ChangedFiles   []string
	Error          error
}

func Status(ctx context.Context, rootDir, workspaceID string) (StatusResult, error) {
	if workspaceID == "" {
		return StatusResult{}, fmt.Errorf("workspace id is required")
	}
	if rootDir == "" {
		return StatusResult{}, fmt.Errorf("root directory is required")
	}
	if err := validateWorkspaceID(ctx, workspaceID); err != nil {
		return StatusResult{}, err
	}

	wsDir := WorkspaceDir(rootDir, workspaceID)
	if exists, err := paths.DirExists(wsDir); err != nil {
		return StatusResult{}, err
	} else if !exists {
		return StatusResult{}, fmt.Errorf("workspace does not exist: %s", wsDir)
	}

	repos, warnings, err := ScanRepos(ctx, wsDir)
	if err != nil {
		return StatusResult{}, err
	}

	result := StatusResult{
		WorkspaceID: workspaceID,
		Warnings:    warnings,
	}
	for _, repo := range repos {
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
		repoStatus.Branch, repoStatus.Upstream, repoStatus.Head, repoStatus.Detached, repoStatus.HeadMissing, repoStatus.Dirty, repoStatus.UntrackedCount, repoStatus.StagedCount, repoStatus.UnstagedCount, repoStatus.UnmergedCount, repoStatus.AheadCount, repoStatus.BehindCount = parseStatusPorcelainV2(statusOut, repoStatus.Branch)
		repoStatus.ChangedFiles = parseChangedFilesPorcelainV2(statusOut)
		result.Repos = append(result.Repos, repoStatus)
	}

	return result, nil
}

func gitStatusPorcelain(ctx context.Context, worktreePath string) (string, error) {
	return gitcmd.StatusPorcelainV2(ctx, worktreePath)
}

func parseStatusPorcelainV2(output, fallbackBranch string) (string, string, string, bool, bool, bool, int, int, int, int, int, int) {
	status := coregitparse.ParseStatusPorcelainV2(output, fallbackBranch)
	return status.Branch, status.Upstream, status.Head, status.Detached, status.HeadMissing, status.Dirty, status.UntrackedCount, status.StagedCount, status.UnstagedCount, status.UnmergedCount, status.AheadCount, status.BehindCount
}

func parseChangedFilesPorcelainV2(output string) []string {
	return coregitparse.ParseChangedFilesPorcelainV2(output)
}
