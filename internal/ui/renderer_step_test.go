package ui

import (
	"bytes"
	"strings"
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

func TestGroupedStepLogger_RendersTreePerStep(t *testing.T) {
	var b bytes.Buffer
	renderer := NewRenderer(&b, DefaultTheme(), false)
	logger := NewGroupedStepLogger(renderer)

	logger.Step("remove workspace test")
	logger.Log("$ git worktree remove --force")
	logger.LogOutput("/tmp/test/repo-a")
	logger.Log("$ git worktree remove --force")
	logger.LogOutput("/tmp/test/repo-b")
	logger.Flush()

	got := b.String()
	wantLines := []string{
		"  • remove workspace test",
		"    ├─ $ git worktree remove --force",
		"    └─ $ git worktree remove --force",
	}
	want := strings.Join(wantLines, "\n") + "\n"
	if got != want {
		t.Fatalf("unexpected output:\n%s\nwant:\n%s", got, want)
	}
}

func TestGroupedStepLogger_RendersGroupedRepoTree(t *testing.T) {
	var b bytes.Buffer
	renderer := NewRenderer(&b, DefaultTheme(), false)
	logger := NewGroupedStepLogger(renderer)

	logger.Step("remove workspace test")
	logger.BeginGroup("repo-a")
	logger.Log("$ git worktree remove --force")
	logger.Log("/tmp/test/repo-a")
	logger.EndGroup()
	logger.BeginGroup("repo-b")
	logger.Log("$ git worktree remove --force")
	logger.Log("/tmp/test/repo-b")
	logger.EndGroup()
	logger.Flush()

	got := b.String()
	wantLines := []string{
		"  • remove workspace test",
		"    ├─ repo-a",
		"    └─ repo-b",
	}
	want := strings.Join(wantLines, "\n") + "\n"
	if got != want {
		t.Fatalf("unexpected output:\n%s\nwant:\n%s", got, want)
	}
}
