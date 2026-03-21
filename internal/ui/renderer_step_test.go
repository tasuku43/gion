package ui

import (
	"bytes"
	"testing"

	"github.com/tasuku43/gion/internal/infra/output"
)

func TestRendererStepLog_FormatsWithConnector(t *testing.T) {
	var b bytes.Buffer
	renderer := NewRenderer(&b, DefaultTheme(), true)

	renderer.StepLog("$ git worktree remove --force")

	got := b.String()
	if got != "  "+output.LogConnector+" $ git worktree remove --force\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}
