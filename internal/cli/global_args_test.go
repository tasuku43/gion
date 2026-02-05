package cli

import (
	"reflect"
	"testing"
)

func TestNormalizeGlobalArgs_MovesNoPromptAfterCommandToFront(t *testing.T) {
	t.Parallel()

	in := []string{"apply", "--no-prompt"}
	want := []string{"--no-prompt", "apply"}
	got, err := normalizeGlobalArgs(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeGlobalArgs() = %#v; want %#v", got, want)
	}
}

func TestNormalizeGlobalArgs_MovesRootAfterCommandToFront(t *testing.T) {
	t.Parallel()

	in := []string{"apply", "--root", "/tmp/gion-root", "--no-prompt"}
	want := []string{"--root", "/tmp/gion-root", "--no-prompt", "apply"}
	got, err := normalizeGlobalArgs(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeGlobalArgs() = %#v; want %#v", got, want)
	}
}

func TestNormalizeGlobalArgs_RespectsDoubleDashTerminator(t *testing.T) {
	t.Parallel()

	in := []string{"apply", "--", "--no-prompt"}
	want := []string{"apply", "--", "--no-prompt"}
	got, err := normalizeGlobalArgs(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeGlobalArgs() = %#v; want %#v", got, want)
	}
}

func TestNormalizeGlobalArgs_RootRequiresValue(t *testing.T) {
	t.Parallel()

	_, err := normalizeGlobalArgs([]string{"apply", "--root", "--no-prompt"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
