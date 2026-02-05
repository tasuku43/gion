package manifestplan

import (
	"context"

	coreplanner "github.com/tasuku43/gion-core/planner"
	"github.com/tasuku43/gion/internal/app/manifestimport"
	"github.com/tasuku43/gion/internal/domain/manifest"
)

type Result struct {
	Desired  manifest.File
	Actual   manifest.File
	Changes  []WorkspaceChange
	Warnings []error
}

func Plan(ctx context.Context, rootDir string) (Result, error) {
	validation, err := manifest.Validate(ctx, rootDir)
	if err != nil {
		return Result{}, err
	}
	if len(validation.Issues) > 0 {
		return Result{}, &manifest.ValidationError{Result: validation}
	}

	desired, err := manifest.Load(rootDir)
	if err != nil {
		return Result{}, err
	}
	actual, warnings, err := manifestimport.Build(ctx, rootDir)
	if err != nil {
		return Result{}, err
	}

	changes := coreplanner.Diff(toInventory(desired), toInventory(actual))

	return Result{
		Desired:  desired,
		Actual:   actual,
		Changes:  changes,
		Warnings: warnings,
	}, nil
}

func toInventory(file manifest.File) coreplanner.Inventory {
	workspaces := make(map[string]coreplanner.Workspace, len(file.Workspaces))
	for id, ws := range file.Workspaces {
		repos := make([]coreplanner.Repo, 0, len(ws.Repos))
		for _, repoEntry := range ws.Repos {
			repos = append(repos, coreplanner.Repo{
				Alias:   repoEntry.Alias,
				RepoKey: repoEntry.RepoKey,
				Branch:  repoEntry.Branch,
			})
		}
		workspaces[id] = coreplanner.Workspace{
			ID:    id,
			Repos: repos,
		}
	}
	return coreplanner.Inventory{Workspaces: workspaces}
}
