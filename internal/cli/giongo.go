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
	"github.com/tasuku43/gwst/internal/domain/workspace"
	"github.com/tasuku43/gwst/internal/infra/paths"
	"github.com/tasuku43/gwst/internal/ui"
)

var isTerminal = isatty.IsTerminal

// RunGiongo is the entrypoint for the giongo binary.
func RunGiongo() error {
	fs := flag.NewFlagSet("giongo", flag.ContinueOnError)
	var rootFlag string
	var printFlag bool
	var helpFlag bool
	var versionFlag bool
	fs.StringVar(&rootFlag, "root", "", "override root")
	fs.BoolVar(&printFlag, "print", false, "print selected path")
	fs.BoolVar(&helpFlag, "help", false, "show help")
	fs.BoolVar(&helpFlag, "h", false, "show help")
	fs.BoolVar(&versionFlag, "version", false, "print version")
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		printGiongoHelp(os.Stdout)
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if versionFlag {
		printVersionFor(os.Stdout, "giongo")
		return nil
	}
	if helpFlag {
		printGiongoHelp(os.Stdout)
		return nil
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unknown argument: %s", fs.Args()[0])
	}
	if !isTerminal(os.Stdin.Fd()) {
		return fmt.Errorf("interactive selection requires a TTY")
	}

	rootDir, err := paths.ResolveRoot(rootFlag)
	if err != nil {
		return err
	}

	ctx := context.Background()
	entries, _, err := workspace.List(rootDir)
	if err != nil {
		return err
	}
	choices, err := buildGiongoWorkspaceChoices(ctx, entries)
	if err != nil {
		return err
	}
	theme := ui.DefaultTheme()
	useColor := isTerminal(os.Stdout.Fd())
	selected, err := ui.PromptWorkspaceRepoSelect("giongo", choices, theme, useColor)
	if err != nil {
		if errors.Is(err, ui.ErrPromptCanceled) {
			return nil
		}
		return err
	}
	if strings.TrimSpace(selected) == "" {
		return nil
	}
	if printFlag {
		fmt.Fprintln(os.Stdout, selected)
	}
	return nil
}

func buildGiongoWorkspaceChoices(ctx context.Context, entries []workspace.Entry) ([]ui.WorkspaceChoice, error) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].WorkspaceID < entries[j].WorkspaceID
	})
	choices := make([]ui.WorkspaceChoice, 0, len(entries))
	for _, entry := range entries {
		repos, _, err := workspace.ScanReposShallow(ctx, entry.WorkspacePath)
		if err != nil {
			return nil, err
		}
		repoChoices := make([]ui.PromptChoice, 0, len(repos))
		for _, repoEntry := range repos {
			name := formatRepoName(repoEntry.Alias, repoEntry.RepoKey)
			label := formatRepoLabel(name, repoEntry.Branch)
			repoKey := displayRepoKey(repoEntry.RepoKey)
			details := make([]string, 0, 2)
			if repoKey != "" {
				details = append(details, fmt.Sprintf("repo: %s", repoKey))
			} else if strings.TrimSpace(repoEntry.RepoSpec) != "" {
				details = append(details, fmt.Sprintf("repo: %s", strings.TrimSpace(repoEntry.RepoSpec)))
			}
			if strings.TrimSpace(repoEntry.Branch) != "" {
				details = append(details, fmt.Sprintf("branch: %s", repoEntry.Branch))
			}
			repoChoices = append(repoChoices, ui.PromptChoice{
				Label:       label,
				Value:       repoEntry.WorktreePath,
				Description: strings.TrimSpace(repoEntry.RepoSpec),
				Details:     details,
			})
		}
		sort.Slice(repoChoices, func(i, j int) bool {
			return repoChoices[i].Label < repoChoices[j].Label
		})
		choices = append(choices, ui.WorkspaceChoice{
			ID:            entry.WorkspaceID,
			WorkspacePath: entry.WorkspacePath,
			Description:   entry.Description,
			Repos:         repoChoices,
		})
	}
	return choices, nil
}
