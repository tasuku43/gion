package manifestls

import (
	"context"
	"fmt"
	"strings"

	coremanifestlsplan "github.com/tasuku43/gion-core/manifestlsplan"
	"github.com/tasuku43/gion/internal/app/manifestplan"
	"github.com/tasuku43/gion/internal/domain/workspace"
)

type DriftStatus = coremanifestlsplan.DriftStatus

const (
	DriftApplied = coremanifestlsplan.DriftApplied
	DriftMissing = coremanifestlsplan.DriftMissing
	DriftDrift   = coremanifestlsplan.DriftDrift
	DriftExtra   = coremanifestlsplan.DriftExtra
)

type Entry struct {
	WorkspaceID  string
	Drift        DriftStatus
	Risk         workspace.WorkspaceStateKind
	Description  string
	HasWorkspace bool
}

type Counts = coremanifestlsplan.Counts

type Result struct {
	ManifestEntries []Entry
	ExtraEntries    []Entry
	Counts          Counts
	Warnings        []error
}

func List(ctx context.Context, rootDir string) (Result, error) {
	plan, err := manifestplan.Plan(ctx, rootDir)
	if err != nil {
		return Result{}, err
	}
	desired := plan.Desired

	fsWorkspaces, fsWarnings, err := workspace.List(rootDir)
	if err != nil {
		return Result{}, err
	}
	warnings := append([]error{}, fsWarnings...)
	desiredWorkspaces := make([]coremanifestlsplan.DesiredWorkspace, 0, len(desired.Workspaces))
	for id := range desired.Workspaces {
		ws := desired.Workspaces[id]
		desiredWorkspaces = append(desiredWorkspaces, coremanifestlsplan.DesiredWorkspace{
			ID:          id,
			Description: strings.TrimSpace(ws.Description),
		})
	}
	filesystemWorkspaceIDs := make([]string, 0, len(fsWorkspaces))
	for _, entry := range fsWorkspaces {
		filesystemWorkspaceIDs = append(filesystemWorkspaceIDs, entry.WorkspaceID)
	}
	layout := coremanifestlsplan.Build(desiredWorkspaces, plan.Changes, filesystemWorkspaceIDs)

	entries := make([]Entry, 0, len(layout.ManifestEntries))
	for _, manifestEntry := range layout.ManifestEntries {
		risk := workspace.WorkspaceStateClean
		if manifestEntry.HasWorkspace {
			state, warn := bestEffortWorkspaceRisk(ctx, rootDir, manifestEntry.WorkspaceID)
			if warn != nil {
				warnings = append(warnings, warn)
			}
			risk = state
		}
		entries = append(entries, Entry{
			WorkspaceID:  manifestEntry.WorkspaceID,
			Drift:        manifestEntry.Drift,
			Risk:         risk,
			Description:  manifestEntry.Description,
			HasWorkspace: manifestEntry.HasWorkspace,
		})
	}

	var extras []Entry
	for _, extra := range layout.ExtraEntries {
		risk, warn := bestEffortWorkspaceRisk(ctx, rootDir, extra.WorkspaceID)
		if warn != nil {
			warnings = append(warnings, warn)
		}
		extras = append(extras, Entry{
			WorkspaceID:  extra.WorkspaceID,
			Drift:        extra.Drift,
			Risk:         risk,
			HasWorkspace: true,
		})
	}

	return Result{
		ManifestEntries: entries,
		ExtraEntries:    extras,
		Counts:          layout.Counts,
		Warnings:        warnings,
	}, nil
}

func bestEffortWorkspaceRisk(ctx context.Context, rootDir, workspaceID string) (workspace.WorkspaceStateKind, error) {
	state, err := workspace.State(ctx, rootDir, workspaceID)
	if err != nil {
		return workspace.WorkspaceStateUnknown, fmt.Errorf("workspace %s state: %w", workspaceID, err)
	}
	return aggregateRiskKind(state.Repos), nil
}

// aggregateRiskKind picks a single workspace risk label from repo risks.
//
// We keep "unknown" as a special-case top priority (can't confidently assert safety).
// When unknown is not present, we use a stable order: dirty > diverged > unpushed.
func aggregateRiskKind(repos []workspace.RepoState) workspace.WorkspaceStateKind {
	hasDirty := false
	hasUnknown := false
	hasDiverged := false
	hasUnpushed := false
	for _, repo := range repos {
		switch repo.Kind {
		case workspace.RepoStateUnknown:
			hasUnknown = true
		case workspace.RepoStateDirty:
			hasDirty = true
		case workspace.RepoStateDiverged:
			hasDiverged = true
		case workspace.RepoStateUnpushed:
			hasUnpushed = true
		}
	}
	switch {
	case hasUnknown:
		return workspace.WorkspaceStateUnknown
	case hasDirty:
		return workspace.WorkspaceStateDirty
	case hasDiverged:
		return workspace.WorkspaceStateDiverged
	case hasUnpushed:
		return workspace.WorkspaceStateUnpushed
	default:
		return workspace.WorkspaceStateClean
	}
}
