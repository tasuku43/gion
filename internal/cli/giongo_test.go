package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
