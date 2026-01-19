package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/tasuku43/gwst/internal/app/manifestimport"
	"github.com/tasuku43/gwst/internal/ui"
)

func runImport(ctx context.Context, rootDir string, args []string, noPrompt bool) error {
	if len(args) == 1 && isHelpArg(args[0]) {
		printImportHelp(os.Stdout)
		return nil
	}
	if len(args) != 0 {
		return fmt.Errorf("usage: gwst import")
	}
	result, err := manifestimport.Import(ctx, rootDir)
	if err != nil {
		return err
	}

	theme := ui.DefaultTheme()
	useColor := isatty.IsTerminal(os.Stdout.Fd())
	renderer := ui.NewRenderer(os.Stdout, theme, useColor)

	var warningLines []string
	for _, warn := range result.Warnings {
		warningLines = append(warningLines, warn.Error())
	}
	if len(warningLines) > 0 {
		renderWarningsSection(renderer, "warnings", warningLines, false)
		renderer.Blank()
	}

	workspaceCount := len(result.Manifest.Workspaces)
	repoCount := 0
	for _, ws := range result.Manifest.Workspaces {
		repoCount += len(ws.Repos)
	}

	renderer.Section("Result")
	renderer.Bullet(fmt.Sprintf("write: %s", result.Path))
	renderer.Bullet(fmt.Sprintf("workspaces: %d", workspaceCount))
	renderer.Bullet(fmt.Sprintf("repos: %d", repoCount))
	return nil
}
