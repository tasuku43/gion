package cli

import (
	"strings"
	"testing"
)

func TestGiongoInitScriptForZsh(t *testing.T) {
	script, err := giongoInitScript("zsh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(script, "giongo() {") {
		t.Fatalf("expected function definition")
	}
	if !strings.Contains(script, "command giongo \"$@\"") {
		t.Fatalf("expected init bypass")
	}
	if !strings.Contains(script, "command giongo --print") {
		t.Fatalf("expected --print wrapper")
	}
}

func TestGiongoInitScriptUnsupportedShell(t *testing.T) {
	_, err := giongoInitScript("fish")
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
}
