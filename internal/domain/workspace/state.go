package workspace

import (
	"context"

	coreworkspacerisk "github.com/tasuku43/gion-core/workspacerisk"
)

type WorkspaceStateKind string

const (
	WorkspaceStateClean    WorkspaceStateKind = "clean"
	WorkspaceStateDirty    WorkspaceStateKind = "dirty"
	WorkspaceStateUnpushed WorkspaceStateKind = "unpushed"
	WorkspaceStateDiverged WorkspaceStateKind = "diverged"
	WorkspaceStateUnknown  WorkspaceStateKind = "unknown"
)

type RepoStateKind string

const (
	RepoStateClean    RepoStateKind = "clean"
	RepoStateDirty    RepoStateKind = "dirty"
	RepoStateUnpushed RepoStateKind = "unpushed"
	RepoStateDiverged RepoStateKind = "diverged"
	RepoStateUnknown  RepoStateKind = "unknown"
)

type WorkspaceState struct {
	WorkspaceID string
	Kind        WorkspaceStateKind
	Repos       []RepoState
	Warnings    []error
}

type RepoState struct {
	Alias          string
	WorktreePath   string
	Upstream       string
	AheadCount     int
	BehindCount    int
	StagedCount    int
	UnstagedCount  int
	UntrackedCount int
	UnmergedCount  int
	Kind           RepoStateKind
	Error          error
}

func State(ctx context.Context, rootDir, workspaceID string) (WorkspaceState, error) {
	status, err := Status(ctx, rootDir, workspaceID)
	if err != nil {
		return WorkspaceState{}, err
	}
	return StateFromStatus(status), nil
}

func StateFromStatus(status StatusResult) WorkspaceState {
	repos := make([]RepoState, 0, len(status.Repos))
	for _, repo := range status.Repos {
		repos = append(repos, repoStateFromStatus(repo))
	}
	return WorkspaceState{
		WorkspaceID: status.WorkspaceID,
		Kind:        aggregateWorkspaceState(repos),
		Repos:       repos,
		Warnings:    status.Warnings,
	}
}

func repoStateFromStatus(repo RepoStatus) RepoState {
	state := RepoState{
		Alias:          repo.Alias,
		WorktreePath:   repo.WorktreePath,
		Upstream:       repo.Upstream,
		AheadCount:     repo.AheadCount,
		BehindCount:    repo.BehindCount,
		StagedCount:    repo.StagedCount,
		UnstagedCount:  repo.UnstagedCount,
		UntrackedCount: repo.UntrackedCount,
		UnmergedCount:  repo.UnmergedCount,
		Error:          repo.Error,
	}
	kind := coreworkspacerisk.ClassifyRepoStatus(coreworkspacerisk.RepoStatus{
		Upstream:    repo.Upstream,
		AheadCount:  repo.AheadCount,
		BehindCount: repo.BehindCount,
		Dirty:       repo.Dirty,
		Detached:    repo.Detached,
		HeadMissing: repo.HeadMissing,
		Error:       repo.Error,
	})
	state.Kind = RepoStateKind(kind)
	return state
}

func aggregateWorkspaceState(repos []RepoState) WorkspaceStateKind {
	coreRepos := make([]coreworkspacerisk.RepoState, 0, len(repos))
	for _, repo := range repos {
		coreRepos = append(coreRepos, coreworkspacerisk.RepoState(repo.Kind))
	}
	switch coreworkspacerisk.AggregateForState(coreRepos) {
	case coreworkspacerisk.WorkspaceRiskDirty:
		return WorkspaceStateDirty
	case coreworkspacerisk.WorkspaceRiskUnknown:
		return WorkspaceStateUnknown
	case coreworkspacerisk.WorkspaceRiskDiverged:
		return WorkspaceStateDiverged
	case coreworkspacerisk.WorkspaceRiskUnpushed:
		return WorkspaceStateUnpushed
	default:
		return WorkspaceStateClean
	}
}

func RequiresRemoveConfirmation(kind WorkspaceStateKind) bool {
	switch kind {
	case WorkspaceStateDirty, WorkspaceStateUnpushed, WorkspaceStateDiverged, WorkspaceStateUnknown:
		return true
	default:
		return false
	}
}
