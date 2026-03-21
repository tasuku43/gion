package workspace

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tasuku43/gion/internal/domain/repo"
)

func TestRemoveDeletesRemoveStateOnSuccess(t *testing.T) {
	t.Setenv("GIT_AUTHOR_NAME", "gion")
	t.Setenv("GIT_AUTHOR_EMAIL", "gion@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "gion")
	t.Setenv("GIT_COMMITTER_EMAIL", "gion@example.com")

	ctx := context.Background()
	rootDir, repoSpec := setupRemoveStateFixture(t, ctx)

	if err := Remove(ctx, rootDir, "WS-1"); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	path, err := removeStatePath(rootDir, "WS-1")
	if err != nil {
		t.Fatalf("removeStatePath() error = %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("remove state should be deleted, stat err = %v", err)
	}
	if _, err := os.Stat(WorkspaceDir(rootDir, "WS-1")); !os.IsNotExist(err) {
		t.Fatalf("workspace should be deleted, stat err = %v", err)
	}
	if _, _, err := repo.Exists(rootDir, repoSpec); err != nil {
		t.Fatalf("repo.Exists() error = %v", err)
	}
}

func TestRemovePersistsPartialRemoveStateWhenWorkspaceDirDeleteFails(t *testing.T) {
	t.Setenv("GIT_AUTHOR_NAME", "gion")
	t.Setenv("GIT_AUTHOR_EMAIL", "gion@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "gion")
	t.Setenv("GIT_COMMITTER_EMAIL", "gion@example.com")

	ctx := context.Background()
	rootDir, _ := setupRemoveStateFixture(t, ctx)

	origRemoveWorkspaceFn := removeWorkspaceFn
	removeWorkspaceFn = func(path string) error {
		return fmt.Errorf("boom: cannot remove %s", path)
	}
	t.Cleanup(func() {
		removeWorkspaceFn = origRemoveWorkspaceFn
	})

	err := Remove(ctx, rootDir, "WS-1")
	if err == nil {
		t.Fatalf("Remove() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "remove workspace dir") {
		t.Fatalf("Remove() error = %v, want workspace dir error", err)
	}

	state, ok, err := LoadRemoveState(rootDir, "WS-1")
	if err != nil {
		t.Fatalf("LoadRemoveState() error = %v", err)
	}
	if !ok {
		t.Fatalf("remove state should exist")
	}
	if got, want := state.RemovedAliases, []string{"repo"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("RemovedAliases = %v, want %v", got, want)
	}
	if len(state.PendingAliases) != 0 {
		t.Fatalf("PendingAliases = %v, want empty", state.PendingAliases)
	}
	if !state.WorkspaceDirPresent {
		t.Fatalf("WorkspaceDirPresent = false, want true")
	}
	if !strings.Contains(state.LastError, "remove workspace dir") {
		t.Fatalf("LastError = %q, want workspace dir message", state.LastError)
	}
	if _, err := os.Stat(WorkspaceDir(rootDir, "WS-1")); err != nil {
		t.Fatalf("workspace should still exist, stat err = %v", err)
	}
}

func setupRemoveStateFixture(t *testing.T, ctx context.Context) (string, string) {
	t.Helper()

	tmp := t.TempDir()
	rootDir := filepath.Join(tmp, "gion")

	remoteBase := filepath.Join(tmp, "remotes")
	remotePath := filepath.Join(remoteBase, "org", "repo.git")
	if err := os.MkdirAll(filepath.Dir(remotePath), 0o755); err != nil {
		t.Fatalf("mkdir remote: %v", err)
	}
	runGitForRemoveState(t, "", "init", "--bare", remotePath)

	seedDir := filepath.Join(tmp, "seed")
	runGitForRemoveState(t, "", "init", seedDir)
	runGitForRemoveState(t, seedDir, "checkout", "-b", "main")
	if err := os.WriteFile(filepath.Join(seedDir, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write seed file: %v", err)
	}
	runGitForRemoveState(t, seedDir, "add", ".")
	runGitForRemoveState(t, seedDir, "commit", "-m", "init")
	runGitForRemoveState(t, seedDir, "remote", "add", "origin", remotePath)
	runGitForRemoveState(t, seedDir, "push", "origin", "main")
	runGitForRemoveState(t, "", "--git-dir", remotePath, "symbolic-ref", "HEAD", "refs/heads/main")

	configPath := filepath.Join(tmp, "gitconfig")
	fileURL := "file://" + filepath.ToSlash(remoteBase) + "/"
	configData := fmt.Sprintf("[url \"%s\"]\n\tinsteadOf = https://example.com/\n", fileURL)
	if err := os.WriteFile(configPath, []byte(configData), 0o644); err != nil {
		t.Fatalf("write gitconfig: %v", err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", configPath)
	t.Setenv("GIT_CONFIG_SYSTEM", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_TERMINAL_PROMPT", "0")

	repoSpec := "https://example.com/org/repo.git"
	if _, err := repo.Get(ctx, rootDir, repoSpec); err != nil {
		t.Fatalf("repo get: %v", err)
	}
	if _, err := New(ctx, rootDir, "WS-1"); err != nil {
		t.Fatalf("workspace new: %v", err)
	}
	if _, err := Add(ctx, rootDir, "WS-1", repoSpec, "", true); err != nil {
		t.Fatalf("workspace add: %v", err)
	}
	return rootDir, repoSpec
}

func runGitForRemoveState(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s failed: %v\nstdout:\n%s\nstderr:\n%s", strings.Join(args, " "), err, stdout.String(), stderr.String())
	}
}
