package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunShellInit_WithCompletion(t *testing.T) {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error: %v", err)
	}
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = stdout })

	err = runShell([]string{"init", "zsh", "--with-completion"})
	if err != nil {
		t.Fatalf("runShell(init) error: %v", err)
	}
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("ReadFrom() error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "GION_SHELL_ACTION_FILE") {
		t.Fatalf("shell init output missing action wrapper: %q", output)
	}
	if !strings.Contains(output, "#compdef gion") {
		t.Fatalf("shell init output missing zsh completion: %q", output)
	}
}

func TestRunShellCompletion(t *testing.T) {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() error: %v", err)
	}
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = stdout })

	err = runShell([]string{"completion", "bash"})
	if err != nil {
		t.Fatalf("runShell(completion) error: %v", err)
	}
	_ = w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("ReadFrom() error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "_gion_completion") {
		t.Fatalf("shell completion output missing bash completion: %q", output)
	}
	if strings.Contains(output, "GION_SHELL_ACTION_FILE") {
		t.Fatalf("shell completion should not include shell wrapper: %q", output)
	}
}

func TestRunShellInit_InvalidWithCompletionValue(t *testing.T) {
	err := runShellInit([]string{"zsh", "--with-completion=maybe"})
	if err == nil {
		t.Fatal("expected error for invalid --with-completion value")
	}
}
