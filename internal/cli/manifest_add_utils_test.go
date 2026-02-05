package cli

import (
	"reflect"
	"testing"
)

func TestNormalizeManifestAddArgs_ReordersFlagsAfterPositionals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "repo_mode_no_prompt_after_workspace_id",
			in:   []string{"--repo", "git@github.com:tasuku43/gionx.git", "MVP-001", "--no-prompt"},
			want: []string{"--repo", "git@github.com:tasuku43/gionx.git", "--no-prompt", "MVP-001"},
		},
		{
			name: "positional_before_repo_flag",
			in:   []string{"MVP-001", "--repo", "git@github.com:tasuku43/gionx.git"},
			want: []string{"--repo", "git@github.com:tasuku43/gionx.git", "MVP-001"},
		},
		{
			name: "missing_repo_value_becomes_empty_assignment",
			in:   []string{"--repo", "--no-prompt", "MVP-001"},
			want: []string{"--repo=", "--no-prompt", "MVP-001"},
		},
		{
			name: "unknown_flag_after_positional_moves_before",
			in:   []string{"--repo", "x", "MVP-001", "--no-promp"},
			want: []string{"--repo", "x", "--no-promp", "MVP-001"},
		},
		{
			name: "double_dash_terminator_keeps_rest_positional",
			in:   []string{"--repo", "x", "MVP-001", "--", "--no-prompt"},
			want: []string{"--repo", "x", "MVP-001", "--", "--no-prompt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeManifestAddArgs(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("normalizeManifestAddArgs() = %#v; want %#v", got, tt.want)
			}
		})
	}
}

func TestNormalizeManifestPresetAddArgs_ReordersRepoAndNoPrompt(t *testing.T) {
	t.Parallel()

	in := []string{"mypreset", "--repo", "a", "--repo", "b", "--no-prompt"}
	want := []string{"--repo", "a", "--repo", "b", "--no-prompt", "mypreset"}
	got := normalizeManifestPresetAddArgs(in)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeManifestPresetAddArgs() = %#v; want %#v", got, want)
	}
}

func TestNormalizeArgsFlagsFirst_ReordersNoPromptForRmStyleArgs(t *testing.T) {
	t.Parallel()

	in := []string{"MVP-001", "--no-prompt"}
	want := []string{"--no-prompt", "MVP-001"}
	got := normalizeArgsFlagsFirst(in, nil)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeArgsFlagsFirst() = %#v; want %#v", got, want)
	}
}
