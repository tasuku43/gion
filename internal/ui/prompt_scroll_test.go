package ui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/tasuku43/gion/internal/infra/output"
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

func TestMultiSelectView_ShowsAssistLines(t *testing.T) {
	model := newMultiSelectModel("title", "repo", []PromptChoice{
		{Label: "example/repo-a", Value: "a"},
		{Label: "example/repo-b", Value: "b"},
	}, DefaultTheme(), false)
	model.selectedValues = []string{"a"}
	model.selected = []PromptChoice{{Label: "example/repo-a", Value: "a"}}
	model.input.SetValue("repo")

	view := model.View()
	if !strings.Contains(view, "filter: repo") {
		t.Fatalf("expected filter assist line, got: %q", view)
	}
	if !strings.Contains(view, "selected: 1/2") {
		t.Fatalf("expected selected summary line, got: %q", view)
	}
	if !strings.Contains(view, "space toggle") {
		t.Fatalf("expected toggle hint line, got: %q", view)
	}
	if !strings.Contains(view, "enter apply") {
		t.Fatalf("expected action hint line, got: %q", view)
	}
}

func TestWorkspaceMultiSelectView_ShowsAssistLines(t *testing.T) {
	model := newWorkspaceMultiSelectModel("gion manifest rm", []WorkspaceChoice{
		{ID: "WS-1", Description: "first"},
		{ID: "WS-2", Description: "second"},
	}, nil, DefaultTheme(), false)
	model.selectedIDs = []string{"WS-1"}
	model.selected = []WorkspaceChoice{{ID: "WS-1", Description: "first"}}
	model.input.SetValue("WS")

	view := model.View()
	if !strings.Contains(view, "filter: WS") {
		t.Fatalf("expected filter assist line, got: %q", view)
	}
	if !strings.Contains(view, "selected: 1/2") {
		t.Fatalf("expected selected summary line, got: %q", view)
	}
	if !strings.Contains(view, "space toggle") {
		t.Fatalf("expected toggle hint line, got: %q", view)
	}
	if !strings.Contains(view, "enter apply") {
		t.Fatalf("expected apply action hint line, got: %q", view)
	}
}

func TestMultiSelectCandidateRows_ShowFocusMarker(t *testing.T) {
	var b strings.Builder
	renderRepoChoiceList(&b, []PromptChoice{
		{Label: "example/repo-a", Value: "a"},
		{Label: "example/repo-b", Value: "b"},
	}, 1, 10, false, DefaultTheme())

	out := b.String()
	if !strings.Contains(out, "> └─ example/repo-b") {
		t.Fatalf("expected focused repo row marker, got: %q", out)
	}
}

func TestWorkspaceMultiSelectCandidateRows_ShowFocusMarker(t *testing.T) {
	var b strings.Builder
	renderWorkspaceChoiceList(&b, []WorkspaceChoice{
		{ID: "WS-1", Description: "first"},
		{ID: "WS-2", Description: "second"},
	}, 0, 10, false, DefaultTheme())

	out := b.String()
	if !strings.Contains(out, "> └─ WS-1 - first") {
		t.Fatalf("expected focused workspace row marker, got: %q", out)
	}
}

func TestRepoMultiSelectRows_ShowSelectionMarker(t *testing.T) {
	var b strings.Builder
	renderRepoMultiSelectChoiceList(&b, []PromptChoice{
		{Label: "example/repo-a", Value: "a"},
		{Label: "example/repo-b", Value: "b"},
	}, 0, 10, map[string]bool{"b": true}, false, DefaultTheme())

	out := b.String()
	if !strings.Contains(out, "○ example/repo-a") || !strings.Contains(out, "● example/repo-b") {
		t.Fatalf("expected selection markers, got: %q", out)
	}
}

func TestWorkspaceMultiSelectRows_ShowSelectionMarker(t *testing.T) {
	var b strings.Builder
	renderWorkspaceMultiSelectChoiceList(&b, []WorkspaceChoice{
		{ID: "WS-1", Description: "first"},
		{ID: "WS-2", Description: "second"},
	}, 0, 10, map[string]bool{"WS-2": true}, false, DefaultTheme())

	out := b.String()
	if !strings.Contains(out, "○ WS-1 - first") || !strings.Contains(out, "● WS-2 - second") {
		t.Fatalf("expected selection markers, got: %q", out)
	}
}

func TestMultiSelectModel_SpaceTogglesSelection(t *testing.T) {
	model := newMultiSelectModel("title", "repo", []PromptChoice{
		{Label: "example/repo-a", Value: "a"},
		{Label: "example/repo-b", Value: "b"},
	}, DefaultTheme(), false)

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeySpace})
	next := updated.(multiSelectModel)
	if len(next.selectedValues) != 1 || next.selectedValues[0] != "a" {
		t.Fatalf("expected first toggle to select a, got: %+v", next.selectedValues)
	}

	updated, _ = next.Update(tea.KeyMsg{Type: tea.KeyUp})
	next = updated.(multiSelectModel)
	updated, _ = next.Update(tea.KeyMsg{Type: tea.KeySpace})
	next = updated.(multiSelectModel)
	if len(next.selectedValues) != 0 {
		t.Fatalf("expected second toggle to clear selection, got: %+v", next.selectedValues)
	}
}

func TestChoiceSelectView_ShowsSingleSelectAssistLines(t *testing.T) {
	model := newChoiceSelectModel("title", "repo", []PromptChoice{
		{Label: "example/repo-a", Value: "a"},
		{Label: "example/repo-b", Value: "b"},
	}, DefaultTheme(), false)
	model.input.SetValue("repo")

	view := model.View()
	if !strings.Contains(view, "filter: repo") {
		t.Fatalf("expected filter assist line, got: %q", view)
	}
	if !strings.Contains(view, "space/enter select") {
		t.Fatalf("expected single-select hint line, got: %q", view)
	}
	if strings.Contains(view, "selected:") {
		t.Fatalf("single-select should not show selected summary, got: %q", view)
	}
}

func TestWorkspaceSelectView_ShowsSingleSelectAssistLines(t *testing.T) {
	model := newWorkspaceSelectModel("title", []WorkspaceChoice{
		{ID: "WS-1", Description: "first"},
		{ID: "WS-2", Description: "second"},
	}, DefaultTheme(), false)
	model.input.SetValue("WS")

	view := model.View()
	if !strings.Contains(view, "filter: WS") {
		t.Fatalf("expected filter assist line, got: %q", view)
	}
	if !strings.Contains(view, "space/enter select") {
		t.Fatalf("expected single-select hint line, got: %q", view)
	}
	if strings.Contains(view, "selected:") {
		t.Fatalf("single-select should not show selected summary, got: %q", view)
	}
}

func TestChoiceSelectModel_SpaceSelectsCurrent(t *testing.T) {
	model := newChoiceSelectModel("title", "repo", []PromptChoice{
		{Label: "example/repo-a", Value: "a"},
		{Label: "example/repo-b", Value: "b"},
	}, DefaultTheme(), false)
	model.cursor = 1

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeySpace})
	next := updated.(choiceSelectModel)
	if next.value != "b" || !next.done {
		t.Fatalf("expected space to select b, got value=%q done=%v", next.value, next.done)
	}
}

func TestWorkspaceSelectModel_SpaceSelectsCurrent(t *testing.T) {
	model := newWorkspaceSelectModel("title", []WorkspaceChoice{
		{ID: "WS-1", Description: "first"},
		{ID: "WS-2", Description: "second"},
	}, DefaultTheme(), false)
	model.cursor = 1

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeySpace})
	next := updated.(workspaceSelectModel)
	if next.workspaceID != "WS-2" {
		t.Fatalf("expected space to select WS-2, got %q", next.workspaceID)
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

func TestBranchInputModel_SeparateInputLineKeepsBranchVisible(t *testing.T) {
	setWrapWidth(60)
	defer setWrapWidth(0)
	setStableLayout(true)
	defer setStableLayout(false)

	model := newBranchInputModel(
		"title",
		[]PromptChoice{{Label: "#96 Refactor dependencies to cli -> app -> domain -> infra", Value: "96"}},
		func(index int, choice PromptChoice) string {
			return fmt.Sprintf("issue #%d (%s)", index+1, choice.Label)
		},
		func(choice PromptChoice) string {
			return "issue/96"
		},
		nil,
		false,
		DefaultTheme(),
		false,
	)
	model.separateInputLine = true

	out := model.ViewWithHeader("repo: git@github.com:tasuku43/gion.git")
	if !strings.Contains(out, "branch:") {
		t.Fatalf("expected output to contain branch line")
	}
	if !strings.Contains(out, "issue/96") {
		t.Fatalf("expected output to contain current input value")
	}
}

func TestConfirmInlineLineModel_IsMultiline(t *testing.T) {
	setWrapWidth(20)
	defer setWrapWidth(0)
	setStableLayout(true)
	defer setStableLayout(false)

	model := newConfirmInlineLineModel("Apply changes? (default: No)", DefaultTheme(), false)
	out := model.View()
	if strings.Count(out, "\n") < 2 {
		t.Fatalf("expected multiline output, got: %q", out)
	}
	if !strings.Contains(out, output.LogConnector) {
		t.Fatalf("expected output to contain connector")
	}
}
