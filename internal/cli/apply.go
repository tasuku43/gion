package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/tasuku43/gwst/internal/app/apply"
	"github.com/tasuku43/gwst/internal/app/manifestplan"
	"github.com/tasuku43/gwst/internal/ui"
)

func runApply(ctx context.Context, rootDir string, args []string, noPrompt bool) error {
	if len(args) == 1 && isHelpArg(args[0]) {
		printApplyHelp(os.Stdout)
		return nil
	}
	if len(args) != 0 {
		return fmt.Errorf("usage: gwst apply")
	}

	plan, err := manifestplan.Plan(ctx, rootDir)
	if err != nil {
		return err
	}

	theme := ui.DefaultTheme()
	useColor := isatty.IsTerminal(os.Stdout.Fd())
	renderer := ui.NewRenderer(os.Stdout, theme, useColor)

	var warningLines []string
	for _, warn := range plan.Warnings {
		warningLines = append(warningLines, warn.Error())
	}
	if len(warningLines) > 0 {
		renderWarningsSection(renderer, "warnings", warningLines, false)
		renderer.Blank()
	}

	renderer.Section("Diff")
	if len(plan.Changes) == 0 {
		renderer.Bullet("no changes")
		return nil
	}
	renderPlanChanges(renderer, plan.Changes)

	destructive := planHasDestructiveChanges(plan)
	if destructive && noPrompt {
		return fmt.Errorf("destructive changes require confirmation")
	}
	if !noPrompt {
		renderer.Blank()
		label := "Apply changes? (default: No)"
		if destructive {
			label = "Apply destructive changes? (default: No)"
		}
		confirm, err := ui.PromptConfirmInline(label, theme, useColor)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	if err := apply.Apply(ctx, rootDir, plan, apply.Options{
		AllowDirty:       destructive,
		AllowStatusError: destructive,
		PrefetchTimeout:  defaultPrefetchTimeout,
	}); err != nil {
		return err
	}

	renderer.Blank()
	renderer.Section("Result")
	renderer.Bullet("applied")
	return nil
}

func planHasDestructiveChanges(plan manifestplan.Result) bool {
	for _, change := range plan.Changes {
		switch change.Kind {
		case manifestplan.WorkspaceRemove:
			return true
		case manifestplan.WorkspaceUpdate:
			if hasDestructiveRepoChange(change.Repos) {
				return true
			}
		}
	}
	return false
}

func hasDestructiveRepoChange(changes []manifestplan.RepoChange) bool {
	for _, change := range changes {
		switch change.Kind {
		case manifestplan.RepoRemove, manifestplan.RepoUpdate:
			return true
		}
	}
	return false
}
