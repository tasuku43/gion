package shellaction

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEmitCDWritesShellActionFile(t *testing.T) {
	actionFile := filepath.Join(t.TempDir(), "action.sh")
	t.Setenv(FileEnv, actionFile)

	if err := EmitCD("/tmp/example dir"); err != nil {
		t.Fatalf("EmitCD() error: %v", err)
	}

	got, err := os.ReadFile(actionFile)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	want := "builtin cd -- '/tmp/example dir'\n"
	if string(got) != want {
		t.Fatalf("shell action = %q, want %q", string(got), want)
	}
}

func TestEmitCDSkipsWhenEnvMissing(t *testing.T) {
	if err := EmitCD("/tmp/example"); err != nil {
		t.Fatalf("EmitCD() error: %v", err)
	}
}
