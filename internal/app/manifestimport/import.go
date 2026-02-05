package manifestimport

import (
	"context"
	"fmt"
	"os"
	"strings"

	coreimportplan "github.com/tasuku43/gion-core/importplan"
	"github.com/tasuku43/gion/internal/domain/manifest"
	"github.com/tasuku43/gion/internal/domain/workspace"
	"github.com/tasuku43/gion/internal/infra/paths"
)

type Result struct {
	Path     string
	Manifest manifest.File
	Warnings []error
}

func Import(ctx context.Context, rootDir string) (Result, error) {
	file, warnings, err := Build(ctx, rootDir)
	if err != nil {
		return Result{}, err
	}
	return Write(rootDir, file, warnings)
}

func Build(ctx context.Context, rootDir string) (manifest.File, []error, error) {
	if strings.TrimSpace(rootDir) == "" {
		return manifest.File{}, nil, fmt.Errorf("root directory is required")
	}
	wsRoot := paths.WorkspacesRoot(rootDir)
	exists, err := paths.DirExists(wsRoot)
	if err != nil {
		return manifest.File{}, nil, err
	}
	if !exists {
		return manifest.File{Version: 1, Workspaces: map[string]manifest.Workspace{}}, nil, nil
	}
	entries, err := os.ReadDir(wsRoot)
	if err != nil {
		return manifest.File{}, nil, err
	}

	file := manifest.File{
		Version:    1,
		Workspaces: map[string]manifest.Workspace{},
	}
	if existing, err := manifest.Load(rootDir); err == nil {
		file.Presets = existing.Presets
	}
	var warnings []error

	var workspaceNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		workspaceNames = append(workspaceNames, entry.Name())
	}
	workspaceIDs := coreimportplan.CollectWorkspaceIDs(workspaceNames)

	snapshots := make([]coreimportplan.WorkspaceSnapshot, 0, len(workspaceIDs))

	for _, wsID := range workspaceIDs {
		wsDir := workspace.WorkspaceDir(rootDir, wsID)
		meta, err := workspace.LoadMetadata(wsDir)
		if err != nil {
			warnings = append(warnings, fmt.Errorf("workspace %s metadata: %w", wsID, err))
		}

		repos, repoWarnings, err := workspace.ScanRepos(ctx, wsDir)
		if err != nil {
			warnings = append(warnings, fmt.Errorf("workspace %s repos: %w", wsID, err))
			continue
		}
		if len(repoWarnings) > 0 {
			for _, warn := range repoWarnings {
				warnings = append(warnings, fmt.Errorf("workspace %s repo: %w", wsID, warn))
			}
		}

		repoEntries := make([]coreimportplan.RepoSnapshot, 0, len(repos))
		for _, repoEntry := range repos {
			repoEntries = append(repoEntries, coreimportplan.RepoSnapshot{
				Alias:   strings.TrimSpace(repoEntry.Alias),
				RepoKey: strings.TrimSpace(repoEntry.RepoKey),
				Branch:  strings.TrimSpace(repoEntry.Branch),
			})
		}

		snapshots = append(snapshots, coreimportplan.WorkspaceSnapshot{
			ID:          wsID,
			Description: strings.TrimSpace(meta.Description),
			Mode:        strings.TrimSpace(meta.Mode),
			PresetName:  strings.TrimSpace(meta.PresetName),
			SourceURL:   strings.TrimSpace(meta.SourceURL),
			BaseBranch:  strings.TrimSpace(meta.BaseBranch),
			Repos:       repoEntries,
		})
	}
	file.Workspaces = toManifestWorkspaces(coreimportplan.BuildInventory(snapshots).Workspaces)

	return file, warnings, nil
}

func Write(rootDir string, file manifest.File, warnings []error) (Result, error) {
	if err := manifest.Save(rootDir, file); err != nil {
		return Result{}, err
	}
	return Result{
		Path:     manifest.Path(rootDir),
		Manifest: file,
		Warnings: warnings,
	}, nil
}

func Path(rootDir string) string {
	return manifest.Path(rootDir)
}

func toManifestWorkspaces(workspaces map[string]coreimportplan.Workspace) map[string]manifest.Workspace {
	converted := make(map[string]manifest.Workspace, len(workspaces))
	for id, ws := range workspaces {
		repos := make([]manifest.Repo, 0, len(ws.Repos))
		for _, repoEntry := range ws.Repos {
			repos = append(repos, manifest.Repo{
				Alias:   repoEntry.Alias,
				RepoKey: repoEntry.RepoKey,
				Branch:  repoEntry.Branch,
				BaseRef: repoEntry.BaseRef,
			})
		}
		converted[id] = manifest.Workspace{
			Description: ws.Description,
			Mode:        ws.Mode,
			PresetName:  ws.PresetName,
			SourceURL:   ws.SourceURL,
			Repos:       repos,
		}
	}
	return converted
}
