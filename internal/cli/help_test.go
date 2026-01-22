package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintCommandHelp_ManifestAliases(t *testing.T) {
	cases := []string{"manifest", "man", "m"}
	for _, cmd := range cases {
		t.Run(cmd, func(t *testing.T) {
			var buf bytes.Buffer
			if ok := printCommandHelp(cmd, &buf); !ok {
				t.Fatalf("expected ok=true")
			}
			out := buf.String()
			if !strings.Contains(out, "Usage: gwst manifest") {
				t.Fatalf("expected manifest usage, got:\n%s", out)
			}
		})
	}
}

func TestPrintCommandHelp_LsRemoved(t *testing.T) {
	var buf bytes.Buffer
	if ok := printCommandHelp("ls", &buf); !ok {
		t.Fatalf("expected ok=true")
	}
	out := buf.String()
	if !strings.Contains(out, "gwst ls is removed") {
		t.Fatalf("expected removed message, got:\n%s", out)
	}
	if !strings.Contains(out, "gwst manifest ls") {
		t.Fatalf("expected suggestion to use manifest ls, got:\n%s", out)
	}
}

func TestPrintCommandHelp_PresetRemoved(t *testing.T) {
	var buf bytes.Buffer
	if ok := printCommandHelp("preset", &buf); !ok {
		t.Fatalf("expected ok=true")
	}
	out := buf.String()
	if !strings.Contains(out, "gwst preset is removed") {
		t.Fatalf("expected removed message, got:\n%s", out)
	}
	if !strings.Contains(out, "gwst manifest preset") {
		t.Fatalf("expected suggestion to use manifest preset, got:\n%s", out)
	}
}
