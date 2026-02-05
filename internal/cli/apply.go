package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	coreapplyplan "github.com/tasuku43/gion-core/applyplan"
	coreplanner "github.com/tasuku43/gion-core/planner"
	"github.com/tasuku43/gion/internal/app/apply"
	"github.com/tasuku43/gion/internal/app/manifestplan"
	"github.com/tasuku43/gion/internal/domain/manifest"
	"github.com/tasuku43/gion/internal/domain/repo"
	"github.com/tasuku43/gion/internal/infra/output"
	"github.com/tasuku43/gion/internal/infra/prefetcher"
	"github.com/tasuku43/gion/internal/ui"
)

func runApply(ctx context.Context, rootDir string, args []string, noPrompt bool) error {
	if len(args) == 1 && isHelpArg(args[0]) {
		printApplyHelp(os.Stdout)
		return nil
	}
	if len(args) != 0 {
		return fmt.Errorf("usage: gion apply")
	}
	plan, err := manifestplan.Plan(ctx, rootDir)
	if err != nil {
		var vErr *manifest.ValidationError
		if errors.As(err, &vErr) {
			theme := ui.DefaultTheme()
			useColor := isatty.IsTerminal(os.Stdout.Fd())
			renderer := ui.NewRenderer(os.Stdout, theme, useColor)
			renderManifestValidationResult(renderer, vErr.Result)
			return err
		}
		return err
	}
	_, err = runApplyInternalWithPlan(ctx, rootDir, nil, noPrompt, plan)
	return err
}

type applyInternalResult struct {
	HadChanges bool
	Confirmed  bool
	Applied    bool
	Canceled   bool
}

func planHasDestructiveChanges(plan manifestplan.Result) bool {
	return coreapplyplan.HasDestructiveChanges(plan.Changes)
}

func runApplyInternalWithPlan(ctx context.Context, rootDir string, renderer *ui.Renderer, noPrompt bool, plan manifestplan.Result) (applyInternalResult, error) {

	theme := ui.DefaultTheme()
	useColor := isatty.IsTerminal(os.Stdout.Fd())
	if renderer == nil {
		renderer = ui.NewRenderer(os.Stdout, theme, useColor)
	}
	output.SetStepLogger(renderer)
	defer output.SetStepLogger(nil)

	var warningLines []string
	for _, warn := range plan.Warnings {
		warningLines = append(warningLines, warn.Error())
	}
	if len(warningLines) > 0 {
		renderWarningsSection(renderer, "warnings", warningLines, false)
		renderer.Blank()
	}

	renderer.Section("Plan")
	if len(plan.Changes) == 0 {
		renderer.Bullet("no changes")
		return applyInternalResult{HadChanges: false, Confirmed: false, Applied: false}, nil
	}
	renderPlanChanges(ctx, rootDir, renderer, plan)

	// Start background fetch while the user reviews the plan.
	// This preserves the "gion manifest add" UX win (fetch overlaps with reading time),
	// while keeping `gion plan` itself side-effect free.
	toPrefetch := repoSpecsForApplyPlan(plan)
	prefetch := prefetcher.New(defaultPrefetchTimeout)
	if _, err := prefetch.StartAll(ctx, rootDir, toPrefetch); err != nil {
		return applyInternalResult{HadChanges: true}, err
	}

	destructive := coreapplyplan.HasDestructiveChanges(plan.Changes)
	if destructive && noPrompt {
		return applyInternalResult{HadChanges: true}, fmt.Errorf("destructive changes require confirmation")
	}
	confirmed := noPrompt
	if !noPrompt {
		renderer.Blank()
		label := "Apply changes? (default: No)"
		if destructive {
			label = "Apply destructive changes? (default: No)"
		}
		confirm, err := ui.PromptConfirmInlinePlan(label, theme, useColor)
		if err != nil {
			if errors.Is(err, ui.ErrPromptCanceled) {
				return applyInternalResult{HadChanges: true, Confirmed: false, Applied: false, Canceled: true}, nil
			}
			return applyInternalResult{HadChanges: true}, err
		}
		confirmed = confirm
		if !confirm {
			return applyInternalResult{HadChanges: true, Confirmed: false, Applied: false}, nil
		}
	}

	renderer.Blank()
	renderer.Section("Apply")
	prefetchOK := true
	if err := prefetch.WaitAll(ctx, toPrefetch); err != nil {
		// ここでのfetch失敗はネットワーク要因が多く、apply自体は継続できることもあるため、
		// 警告を出しつつ続行する。
		renderer.BulletWarn(fmt.Sprintf("prefetch failed (continuing): %v", err))
		prefetchOK = false
	}
	if err := apply.Apply(ctx, rootDir, plan, apply.Options{
		AllowDirty:       destructive,
		AllowStatusError: destructive,
		PrefetchTimeout:  defaultPrefetchTimeout,
		PrefetchOK:       prefetchOK,
		Step:             output.Step,
	}); err != nil {
		return applyInternalResult{HadChanges: true, Confirmed: confirmed, Applied: false}, err
	}
	if err := rebuildManifest(ctx, rootDir); err != nil {
		return applyInternalResult{HadChanges: true, Confirmed: confirmed, Applied: false}, err
	}

	renderer.Blank()
	renderer.Section("Result")
	adds, updates, removes := coreapplyplan.CountWorkspaceChanges(plan.Changes)
	renderer.BulletSuccess(fmt.Sprintf("applied: add=%d update=%d remove=%d", adds, updates, removes))
	renderer.Bullet(fmt.Sprintf("%s rewritten", manifest.FileName))
	return applyInternalResult{HadChanges: true, Confirmed: confirmed, Applied: true}, nil
}

func repoSpecsForApplyPlan(plan manifestplan.Result) []string {
	repoKeys := coreapplyplan.CollectPrefetchRepoKeys(plan.Changes, toPlannerInventory(plan.Desired))
	specs := make([]string, 0, len(repoKeys))
	for _, repoKey := range repoKeys {
		spec := strings.TrimSpace(repo.SpecFromKey(repoKey))
		if spec == "" {
			continue
		}
		specs = append(specs, spec)
	}
	return specs
}

func toPlannerInventory(file manifest.File) coreplanner.Inventory {
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
