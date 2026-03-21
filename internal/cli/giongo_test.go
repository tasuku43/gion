package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tasuku43/gion/internal/domain/workspace"
	"github.com/tasuku43/gion/internal/infra/shellaction"
)

func TestRunGiongoRequiresTTY(t *testing.T) {
	originalArgs := os.Args
	originalIsTerminal := isTerminal
	defer func() {
		os.Args = originalArgs
		isTerminal = originalIsTerminal
	}()

	os.Args = []string{"giongo"}
	isTerminal = func(fd uintptr) bool { return false }

	err := RunGiongo()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "TTY") {
		t.Fatalf("expected TTY error, got %q", err.Error())
	}
}

func TestRunGiongoPrintUsesStderr(t *testing.T) {
	originalArgs := os.Args
	originalIsTerminal := isTerminal
	defer func() {
		os.Args = originalArgs
		isTerminal = originalIsTerminal
	}()

	os.Args = []string{"giongo", "--print"}
	isTerminal = func(fd uintptr) bool { return false }

	err := RunGiongo()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "TTY") {
		t.Fatalf("expected TTY error, got %q", err.Error())
	}
}

func TestFinalizeGiongoSelection_PrintsWhenRequested(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = origStdout }()

	if err := finalizeGiongoSelection("/tmp/example", true); err != nil {
		t.Fatalf("finalizeGiongoSelection() error: %v", err)
	}
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("ReadFrom() error: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "/tmp/example" {
		t.Fatalf("stdout = %q", buf.String())
	}
}

func TestFinalizeGiongoSelection_EmitsShellAction(t *testing.T) {
	actionFile := filepath.Join(t.TempDir(), "action.sh")
	t.Setenv(shellaction.FileEnv, actionFile)

	if err := finalizeGiongoSelection("/tmp/example", false); err != nil {
		t.Fatalf("finalizeGiongoSelection() error: %v", err)
	}

	data, err := os.ReadFile(actionFile)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "builtin cd -- '/tmp/example'\n" {
		t.Fatalf("action file = %q", string(data))
	}
}

func TestBuildGiongoWorkspaceChoices_OmitsDuplicateBranchDetail(t *testing.T) {
	root := t.TempDir()
	wsPath := filepath.Join(root, "WS-1")
	repoPath := filepath.Join(wsPath, "gion")
	if err := os.MkdirAll(repoPath, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = repoPath
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
		}
	}
	run("init", "-b", "feature/test")
	run("remote", "add", "origin", "git@github.com:tasuku43/gion.git")

	choices, err := buildGiongoWorkspaceChoices(context.Background(), []workspace.Entry{{
		WorkspaceID:   "WS-1",
		WorkspacePath: wsPath,
	}})
	if err != nil {
		t.Fatalf("buildGiongoWorkspaceChoices() error: %v", err)
	}
	if len(choices) != 1 || len(choices[0].Repos) != 1 {
		t.Fatalf("unexpected choices: %+v", choices)
	}
	repoChoice := choices[0].Repos[0]
	if repoChoice.Label != "gion (branch: feature/test)" {
		t.Fatalf("repo label = %q", repoChoice.Label)
	}
	if len(repoChoice.Details) != 1 || repoChoice.Details[0] != "repo: github.com/tasuku43/gion" {
		t.Fatalf("repo details = %#v", repoChoice.Details)
	}
}
