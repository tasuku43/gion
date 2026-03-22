package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tasuku43/gion/internal/ui"
)

func TestManifestMutation_PreludeAndInfoHaveBlankLineBetween(t *testing.T) {
	var b bytes.Buffer
	r := ui.NewRenderer(&b, ui.DefaultTheme(), false)

	r.Section("Inputs")
	r.Prompt("mode: preset")
	r.Prompt("preset: gion")
	r.Prompt("workspace id: test")
	r.Blank()
	r.Section("Info")
	r.Bullet("manifest: updated gion.yaml")

	out := b.String()
	if !strings.Contains(out, "workspace id: test\n\nInfo\n") {
		t.Fatalf("expected a blank line before Info section, got: %q", out)
	}
}
