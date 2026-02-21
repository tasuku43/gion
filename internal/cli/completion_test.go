package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintBashCompletion(t *testing.T) {
	var buf bytes.Buffer
	printBashCompletion(&buf)
	output := buf.String()

	if !strings.Contains(output, "_gion_completion") {
		t.Error("missing _gion_completion function")
	}
	if !strings.Contains(output, "complete -F _gion_completion gion") {
		t.Error("missing complete command")
	}
}

func TestPrintZshCompletion(t *testing.T) {
	var buf bytes.Buffer
	printZshCompletion(&buf)
	output := buf.String()

	if !strings.Contains(output, "#compdef gion") {
		t.Error("missing #compdef directive")
	}
	if !strings.Contains(output, "_gion") {
		t.Error("missing _gion function")
	}
}

func TestRunCompletion_Default(t *testing.T) {
	// Default shell is bash
	err := runCompletion(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunCompletion_Invalid(t *testing.T) {
	err := runCompletion([]string{"fish"})
	if err == nil {
		t.Error("expected error for unsupported shell")
	}
	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestSupportedShells(t *testing.T) {
	expected := "bash, zsh"
	if SupportedShells != expected {
		t.Errorf("SupportedShells = %q, want %q", SupportedShells, expected)
	}
}
