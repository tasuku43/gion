package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/tasuku43/gwst/internal/app/manifestplan"
	"github.com/tasuku43/gwst/internal/domain/manifest"
	"github.com/tasuku43/gwst/internal/domain/workspace"
	"github.com/tasuku43/gwst/internal/ui"
)

func runManifestRm(ctx context.Context, rootDir string, args []string, globalNoPrompt bool) error {
	rmFlags := flag.NewFlagSet("manifest rm", flag.ContinueOnError)
	var noApply bool
	var noPromptFlag bool
	var helpFlag bool
	rmFlags.BoolVar(&noApply, "no-apply", false, "do not run gwst apply")
	rmFlags.BoolVar(&noPromptFlag, "no-prompt", false, "disable interactive prompt")
	rmFlags.BoolVar(&helpFlag, "help", false, "show help")
	rmFlags.BoolVar(&helpFlag, "h", false, "show help")
	rmFlags.SetOutput(os.Stdout)
	rmFlags.Usage = func() {
		printManifestRmHelp(os.Stdout)
	}
	if err := rmFlags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if helpFlag {
		printManifestRmHelp(os.Stdout)
		return nil
	}

	noPrompt := globalNoPrompt || noPromptFlag

	desired, err := manifest.Load(rootDir)
	if err != nil {
		return err
	}

	manifestPath := manifest.Path(rootDir)
	originalBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("read gwst.yaml: %w", err)
	}

	selectedIDs := uniqueNonEmptyStrings(rmFlags.Args())
	if len(selectedIDs) == 0 {
		if noPrompt {
			return fmt.Errorf("workspace id is required when --no-prompt is set")
		}
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			return fmt.Errorf("interactive workspace selection requires a TTY")
		}
		theme := ui.DefaultTheme()
		useColor := isatty.IsTerminal(os.Stdout.Fd())
		choices := buildManifestRmWorkspaceChoices(ctx, rootDir, desired)
		selected, err := ui.PromptWorkspaceMultiSelectWithBlocked("gwst manifest rm", choices, nil, theme, useColor)
		if err != nil {
			return err
		}
		if len(selected) == 0 {
			return nil
		}
		selectedIDs = uniqueNonEmptyStrings(selected)
	}
	if len(selectedIDs) == 0 {
		return nil
	}

	var missing []string
	for _, id := range selectedIDs {
		if _, ok := desired.Workspaces[id]; !ok {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("workspace(s) not found in gwst.yaml: %s", strings.Join(missing, ", "))
	}

	updated := desired
	if updated.Workspaces == nil {
		updated.Workspaces = map[string]manifest.Workspace{}
	}
	for _, id := range selectedIDs {
		delete(updated.Workspaces, id)
	}

	inputs := func(r *ui.Renderer) {
		r.Section("Inputs")
		for _, id := range selectedIDs {
			r.Bullet(fmt.Sprintf("workspace: %s", id))
		}
	}

	return applyManifestMutation(ctx, rootDir, updated, manifestMutationOptions{
		NoApply:       noApply,
		NoPrompt:      noPrompt,
		OriginalBytes: originalBytes,
		Hooks: manifestMutationHooks{
			ShowPrelude: inputs,
			RenderNoApply: func(r *ui.Renderer) {
				r.Section("Result")
				r.Bullet(fmt.Sprintf("updated gwst.yaml (removed %d workspace(s))", len(selectedIDs)))
				r.Blank()
				r.Section("Suggestion")
				r.Bullet("gwst apply")
			},
			RenderNoChanges: func(r *ui.Renderer) {
				r.Section("Result")
				r.Bullet(fmt.Sprintf("updated gwst.yaml (removed %d workspace(s))", len(selectedIDs)))
				r.Bullet("no changes")
			},
			RenderInfoBeforeApply: func(r *ui.Renderer, _ manifestplan.Result, _ bool) {
				r.Section("Info")
				r.Bullet(fmt.Sprintf("manifest: updated gwst.yaml (removed %d workspace(s))", len(selectedIDs)))
				r.Bullet("apply: reconciling entire root (destructive removals require confirmation)")
			},
		},
	})
}

func buildManifestRmWorkspaceChoices(ctx context.Context, rootDir string, desired manifest.File) []ui.WorkspaceChoice {
	var ids []string
	for id := range desired.Workspaces {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	choices := make([]ui.WorkspaceChoice, 0, len(ids))
	for _, id := range ids {
		ws := desired.Workspaces[id]
		kind := bestEffortWorkspaceRiskKind(ctx, rootDir, id)
		warn, strong := manifestRmRiskWarning(kind)
		choices = append(choices, ui.WorkspaceChoice{
			ID:            id,
			Description:   strings.TrimSpace(ws.Description),
			Warning:       warn,
			WarningStrong: strong,
		})
	}
	return choices
}

func bestEffortWorkspaceRiskKind(ctx context.Context, rootDir, workspaceID string) workspace.WorkspaceStateKind {
	state, err := workspace.State(ctx, rootDir, workspaceID)
	if err != nil {
		return workspace.WorkspaceStateUnknown
	}
	hasDirty := false
	hasUnknown := false
	hasDiverged := false
	hasUnpushed := false
	for _, repo := range state.Repos {
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

func manifestRmRiskWarning(kind workspace.WorkspaceStateKind) (warning string, strong bool) {
	switch kind {
	case workspace.WorkspaceStateUnknown:
		return "unknown", true
	case workspace.WorkspaceStateDirty:
		return "dirty", true
	case workspace.WorkspaceStateDiverged:
		return "diverged", false
	case workspace.WorkspaceStateUnpushed:
		return "unpushed", false
	default:
		return "", false
	}
}

func uniqueNonEmptyStrings(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
