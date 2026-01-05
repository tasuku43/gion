package gc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tasuku43/gws/internal/workspace"
)

type Candidate struct {
	WorkspaceID   string
	WorkspacePath string
	LastUsedAt    string
	Reason        string
}

type Result struct {
	Candidates []Candidate
	Warnings   []error
}

type Options struct {
	OlderThan time.Duration
}

func DryRun(ctx context.Context, rootDir string, opts Options, now time.Time) (Result, error) {
	if rootDir == "" {
		return Result{}, fmt.Errorf("root directory is required")
	}
	_ = ctx

	entries, warnings, err := workspace.List(rootDir)
	if err != nil {
		return Result{}, err
	}

	var candidates []Candidate
	for _, entry := range entries {
		if entry.Manifest == nil {
			continue
		}
		if entry.Manifest.Policy.Pinned {
			continue
		}

		lastUsedAt, ok := parseRFC3339(entry.Manifest.LastUsedAt)
		if !ok {
			warnings = append(warnings, fmt.Errorf("workspace %s: invalid last_used_at: %q", entry.WorkspaceID, entry.Manifest.LastUsedAt))
			continue
		}

		if opts.OlderThan > 0 {
			if now.Sub(lastUsedAt) < opts.OlderThan {
				continue
			}
			candidates = append(candidates, Candidate{
				WorkspaceID:   entry.WorkspaceID,
				WorkspacePath: entry.WorkspacePath,
				LastUsedAt:    entry.Manifest.LastUsedAt,
				Reason:        fmt.Sprintf("older_than:%s", opts.OlderThan),
			})
			continue
		}

		ttlDays := entry.Manifest.Policy.TTLDays
		if ttlDays <= 0 {
			continue
		}
		if now.Sub(lastUsedAt) < time.Duration(ttlDays)*24*time.Hour {
			continue
		}
		candidates = append(candidates, Candidate{
			WorkspaceID:   entry.WorkspaceID,
			WorkspacePath: entry.WorkspacePath,
			LastUsedAt:    entry.Manifest.LastUsedAt,
			Reason:        fmt.Sprintf("ttl_days:%d", ttlDays),
		})
	}

	return Result{Candidates: candidates, Warnings: warnings}, nil
}

func Run(ctx context.Context, rootDir string, opts Options, now time.Time) (Result, error) {
	result, err := DryRun(ctx, rootDir, opts, now)
	if err != nil {
		return Result{}, err
	}
	for _, candidate := range result.Candidates {
		if err := workspace.Remove(ctx, rootDir, candidate.WorkspaceID); err != nil {
			return Result{}, fmt.Errorf("remove workspace %s: %w", candidate.WorkspaceID, err)
		}
		if err := ensureWorkspaceRemoved(candidate.WorkspacePath); err != nil {
			return Result{}, fmt.Errorf("remove workspace %s: %w", candidate.WorkspaceID, err)
		}
	}
	return result, nil
}

func ensureWorkspaceRemoved(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("workspace path still exists: %s", path)
	}
	return nil
}

func parseRFC3339(value string) (time.Time, bool) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, false
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

func WorkspaceRoot(rootDir string) string {
	return filepath.Join(rootDir, "ws")
}
