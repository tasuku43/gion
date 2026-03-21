package apply

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	coreapplyplan "github.com/tasuku43/gion-core/applyplan"
	"github.com/tasuku43/gion/internal/app/manifestplan"
	"github.com/tasuku43/gion/internal/domain/workspace"
)

func TestUpdateBaseBranchCandidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		candidate     string
		mixed         bool
		input         string
		wantCandidate string
		wantMixed     bool
	}{
		{
			name:          "empty input keeps empty",
			candidate:     "",
			mixed:         false,
			input:         "",
			wantCandidate: "",
			wantMixed:     false,
		},
		{
			name:          "first base branch sets candidate",
			candidate:     "",
			mixed:         false,
			input:         "origin/main",
			wantCandidate: "origin/main",
			wantMixed:     false,
		},
		{
			name:          "same base branch keeps candidate",
			candidate:     "origin/main",
			mixed:         false,
			input:         "origin/main",
			wantCandidate: "origin/main",
			wantMixed:     false,
		},
		{
			name:          "different base branch marks mixed",
			candidate:     "origin/main",
			mixed:         false,
			input:         "origin/master",
			wantCandidate: "origin/main",
			wantMixed:     true,
		},
		{
			name:          "once mixed stays mixed",
			candidate:     "origin/main",
			mixed:         true,
			input:         "origin/main",
			wantCandidate: "origin/main",
			wantMixed:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCandidate, gotMixed := coreapplyplan.UpdateBaseBranchCandidate(tt.candidate, tt.mixed, tt.input)
			if gotCandidate != tt.wantCandidate {
				t.Fatalf("candidate: got %q, want %q", gotCandidate, tt.wantCandidate)
			}
			if gotMixed != tt.wantMixed {
				t.Fatalf("mixed: got %v, want %v", gotMixed, tt.wantMixed)
			}
		})
	}
}

func TestApplyFailsFastWhenWorkspaceRemoveStateExists(t *testing.T) {
	t.Parallel()

	rootDir := t.TempDir()
	if err := workspace.SaveRemoveState(rootDir, workspace.NewRemoveState("WS-1", []string{"repo-a", "repo-b"}, true)); err != nil {
		t.Fatalf("SaveRemoveState() error = %v", err)
	}

	err := Apply(context.Background(), rootDir, manifestplan.Result{
		Changes: []manifestplan.WorkspaceChange{
			{
				Kind:        manifestplan.WorkspaceRemove,
				WorkspaceID: "WS-1",
			},
		},
	}, Options{})
	if err == nil {
		t.Fatalf("Apply() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "workspace removal is incomplete for WS-1") {
		t.Fatalf("Apply() error = %v, want incomplete removal message", err)
	}
	wantPath := filepath.Join(rootDir, workspace.MetadataDirName, "state", "operations", "workspace-remove", "WS-1.json")
	if !strings.Contains(err.Error(), wantPath) {
		t.Fatalf("Apply() error = %v, want path %q", err, wantPath)
	}
}
