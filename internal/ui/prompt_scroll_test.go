package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/tasuku43/gwst/internal/infra/output"
)

func countTerminalLines(text string) int {
	trimmed := strings.TrimRight(text, "\n")
	if trimmed == "" {
		return 0
	}
	return strings.Count(trimmed, "\n") + 1
}

func TestWorkspaceRepoSelectView_DoesNotExceedHeight(t *testing.T) {
	setWrapWidth(30)
	defer setWrapWidth(0)

	var workspaces []WorkspaceChoice
	for i := 0; i < 8; i++ {
		var repos []PromptChoice
		for j := 0; j < 8; j++ {
			repos = append(repos, PromptChoice{
				Label:   "repo",
				Value:   "/ws/repo",
				Details: []string{"repo: github.com/tasuku43/gion", "branch: issue/999"},
			})
		}
		workspaces = append(workspaces, WorkspaceChoice{
			ID:          "TASUKU43-GION-ISSUE-999",
			Description: "Refactor dependencies to cli -> app -> domain -> infra",
			Repos:       repos,
		})
	}

	model := newWorkspaceRepoSelectModel("giongo", workspaces, DefaultTheme(), false)
	model.height = 10
	if len(model.selections) > 0 {
		model.cursor = len(model.selections) - 1
	}

	view := model.View()
	if got := countTerminalLines(view); got > model.height {
		t.Fatalf("expected view lines <= %d, got %d", model.height, got)
	}
}

func TestWorkspaceMultiSelectView_DoesNotExceedHeight(t *testing.T) {
	setWrapWidth(30)
	defer setWrapWidth(0)

	var workspaces []WorkspaceChoice
	for i := 0; i < 50; i++ {
		workspaces = append(workspaces, WorkspaceChoice{
			ID:          "TASUKU43-GION-ISSUE-999",
			Description: "Minimal hook/setup automation for workspaces",
			Warning:     "dirty",
		})
	}

	model := newWorkspaceMultiSelectModel("gion manifest rm", workspaces, nil, DefaultTheme(), false)
	model.height = 10
	view := model.View()
	if got := countTerminalLines(view); got > model.height {
		t.Fatalf("expected view lines <= %d, got %d", model.height, got)
	}
}

func TestMultiSelectView_DoesNotExceedHeight(t *testing.T) {
	setWrapWidth(20)
	defer setWrapWidth(0)

	var choices []PromptChoice
	for i := 0; i < 50; i++ {
		choices = append(choices, PromptChoice{Label: "github.com/tasuku43/gion", Value: "value", Description: "very long description for wrapping"})
	}
	model := newMultiSelectModel("title", "repo", choices, DefaultTheme(), false)
	model.height = 10
	view := model.View()
	if got := countTerminalLines(view); got > model.height {
		t.Fatalf("expected view lines <= %d, got %d", model.height, got)
	}
}

func TestStableLayout_TruncatesLinesWithDots(t *testing.T) {
	setWrapWidth(20)
	defer setWrapWidth(0)
	setStableLayout(true)
	defer setStableLayout(false)

	f := NewFrame(DefaultTheme(), false)
	f.SetInputsRaw(output.Indent + output.StepPrefix + " repo: github.com/tasuku43/gion")
	out := f.Render()

	if !strings.Contains(out, "...") {
		t.Fatalf("expected output to contain truncation tail")
	}

	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Section headers are not prefixed/wrapped; they can exceed in extremely narrow terminals.
		if strings.TrimSpace(line) == "Inputs" {
			continue
		}
		if w := ansi.StringWidth(line); w > 20 {
			t.Fatalf("expected line width <= 20, got %d (%q)", w, line)
		}
	}
}
