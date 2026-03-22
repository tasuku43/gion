package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCreateFlow_PresetBranchInput_UsesSeparateInputLine(t *testing.T) {
	m := createFlowModel{
		title:          "gion manifest add",
		mode:           "preset",
		theme:          DefaultTheme(),
		useColor:       false,
		validateBranch: func(string) error { return nil },
	}
	m.presetModel = newInputsModelWithLabel(m.title, []string{"app"}, "app", "PROJ-123", "preset", nil, m.theme, m.useColor)
	m.presetRepos = []string{"git@github.com:org/repo.git"}
	m.beginDescriptionStage()

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := next.(createFlowModel)

	if got.stage != createStagePresetBranch {
		t.Fatalf("expected stage %v, got %v", createStagePresetBranch, got.stage)
	}
	if !got.branchModel.separateInputLine {
		t.Fatalf("expected separateInputLine=true for preset branch input")
	}
}

func TestCreateFlow_PresetSelect_PropagatesConfirmTimerAndAdvances(t *testing.T) {
	m := newCreateFlowModel(
		"gion manifest add",
		[]string{"gion"},
		nil,
		nil,
		nil,
		"",
		"",
		nil,
		nil,
		nil,
		nil,
		func(string) ([]string, error) { return []string{"git@github.com:tasuku43/gion.git"}, nil },
		nil,
		nil,
		nil,
		DefaultTheme(),
		false,
		"preset",
		"",
	)

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	next := updated.(createFlowModel)
	if !next.presetModel.confirming {
		t.Fatalf("expected preset selector to enter confirming state")
	}
	if cmd == nil {
		t.Fatalf("expected preset selector to return confirm timer command")
	}

	updated, cmd = next.Update(singleSelectConfirmDoneMsg{})
	next = updated.(createFlowModel)
	if next.stage != createStagePreset {
		t.Fatalf("expected create flow to remain on preset stage until workspace id is entered, got %v", next.stage)
	}
	if next.presetModel.stage != stageWorkspace {
		t.Fatalf("expected preset selector to advance to workspace input stage, got %v", next.presetModel.stage)
	}
	if cmd != nil {
		t.Fatalf("expected no quit command after advancing to workspace input stage")
	}
}
