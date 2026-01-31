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
