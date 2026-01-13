package cli

import "testing"

func TestFormatReviewWorkspaceID(t *testing.T) {
	cases := []struct {
		name   string
		owner  string
		repo   string
		number int
		want   string
	}{
		{
			name:   "basic",
			owner:  "owner",
			repo:   "repo",
			number: 123,
			want:   "OWNER-REPO-REVIEW-PR-123",
		},
		{
			name:   "trim and uppercase",
			owner:  "  MiXeD ",
			repo:   "  rEpO-name ",
			number: 9,
			want:   "MIXED-REPO-NAME-REVIEW-PR-9",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatReviewWorkspaceID(tc.owner, tc.repo, tc.number); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestFormatIssueWorkspaceID(t *testing.T) {
	cases := []struct {
		name   string
		owner  string
		repo   string
		number int
		want   string
	}{
		{
			name:   "basic",
			owner:  "owner",
			repo:   "repo",
			number: 456,
			want:   "OWNER-REPO-ISSUE-456",
		},
		{
			name:   "trim and uppercase",
			owner:  "  MiXeD ",
			repo:   "  rEpO-name ",
			number: 1,
			want:   "MIXED-REPO-NAME-ISSUE-1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatIssueWorkspaceID(tc.owner, tc.repo, tc.number); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
