package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/tasuku43/gwst/internal/app/manifestimport"
	"github.com/tasuku43/gwst/internal/domain/manifest"
	"github.com/tasuku43/gwst/internal/ui"
)

func runPlan(ctx context.Context, rootDir string, args []string) error {
	if len(args) == 1 && isHelpArg(args[0]) {
		printPlanHelp(os.Stdout)
		return nil
	}
	if len(args) != 0 {
		return fmt.Errorf("usage: gwst plan")
	}

	theme := ui.DefaultTheme()
	useColor := isatty.IsTerminal(os.Stdout.Fd())
	renderer := ui.NewRenderer(os.Stdout, theme, useColor)

	desired, err := manifest.Load(rootDir)
	if err != nil {
		return err
	}
	actual, warnings, err := manifestimport.Build(ctx, rootDir)
	if err != nil {
		return err
	}
	actualBytes, err := manifest.Marshal(actual)
	if err != nil {
		return err
	}
	desiredBytes, err := manifest.Marshal(desired)
	if err != nil {
		return err
	}
	diffLines, err := buildUnifiedDiffLines(actualBytes, desiredBytes)
	if err != nil {
		return err
	}

	var warningLines []string
	for _, warn := range warnings {
		warningLines = append(warningLines, warn.Error())
	}
	if len(warningLines) > 0 {
		renderWarningsSection(renderer, "warnings", warningLines, false)
		renderer.Blank()
	}

	renderer.Section("Diff")
	if len(diffLines) == 0 {
		renderer.Bullet("no changes")
		return nil
	}
	renderDiffLines(renderer, diffLines)
	return nil
}
