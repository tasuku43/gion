package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	coregitparse "github.com/tasuku43/gion-core/gitparse"
	coregitref "github.com/tasuku43/gion-core/gitref"
	corerepostore "github.com/tasuku43/gion-core/repostore"
	"github.com/tasuku43/gion/internal/infra/gitcmd"
	"github.com/tasuku43/gion/internal/infra/paths"
)

type Store struct {
	RepoKey   string
	StorePath string
	RemoteURL string
}

func Get(ctx context.Context, rootDir string, repo string) (Store, error) {
	spec, remoteURL, err := Normalize(repo)
	if err != nil {
		return Store{}, err
	}

	storePath := storePathForSpec(rootDir, spec)

	exists, err := paths.DirExists(storePath)
	if err != nil {
		return Store{}, err
	}

	if !exists {
		if err := os.MkdirAll(filepath.Dir(storePath), 0o750); err != nil {
			return Store{}, fmt.Errorf("create repo store dir: %w", err)
		}
		gitcmd.Logf("git clone --bare %s %s", remoteURL, storePath)
		if _, err := gitcmd.Run(ctx, []string{"clone", "--bare", remoteURL, storePath}, gitcmd.Options{}); err != nil {
			return Store{}, err
		}
	}

	if err := normalizeStore(ctx, storePath, spec.RepoKey, false); err != nil {
		return Store{}, err
	}

	return Store{
		RepoKey:   spec.RepoKey,
		StorePath: storePath,
		RemoteURL: remoteURL,
	}, nil
}

func Open(ctx context.Context, rootDir string, repo string, fetch bool) (Store, error) {
	spec, remoteURL, err := Normalize(repo)
	if err != nil {
		return Store{}, err
	}

	storePath := storePathForSpec(rootDir, spec)

	exists, err := paths.DirExists(storePath)
	if err != nil {
		return Store{}, err
	}
	if !exists {
		return Store{}, fmt.Errorf("repo store not found, run: gion repo get %s", repo)
	}

	if err := normalizeStore(ctx, storePath, spec.RepoKey, fetch); err != nil {
		return Store{}, err
	}

	return Store{
		RepoKey:   spec.RepoKey,
		StorePath: storePath,
		RemoteURL: remoteURL,
	}, nil
}

func Prefetch(ctx context.Context, rootDir string, repo string) error {
	spec, _, err := Normalize(repo)
	if err != nil {
		return err
	}

	storePath := storePathForSpec(rootDir, spec)

	exists, err := paths.DirExists(storePath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("repo store not found, run: gion repo get %s", repo)
	}

	_, err = ensureDefaultBranch(ctx, storePath, true, false)
	return err
}

func Exists(rootDir, repo string) (string, bool, error) {
	spec, _, err := Normalize(repo)
	if err != nil {
		return "", false, err
	}
	storePath := storePathForSpec(rootDir, spec)
	exists, err := paths.DirExists(storePath)
	if err != nil {
		return "", false, err
	}
	return storePath, exists, nil
}

func storePathForSpec(rootDir string, spec Spec) string {
	return StorePath(rootDir, spec)
}

// (moved to paths.go)

func normalizeStore(ctx context.Context, storePath, display string, fetch bool) error {
	defaultBranch, err := ensureDefaultBranch(ctx, storePath, fetch, true)
	if err != nil {
		return err
	}
	return pruneLocalHeads(ctx, storePath, defaultBranch)
}

func ensureDefaultBranch(ctx context.Context, storePath string, fetch bool, log bool) (string, error) {
	if _, err := gitcmd.Run(ctx, []string{"config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*"}, gitcmd.Options{Dir: storePath}); err != nil {
		return "", err
	}
	defaultBranch, _ := localDefaultBranch(ctx, storePath)

	remoteChecked := false
	grace := fetchGraceDuration()
	if !recentlyFetched(storePath, grace) || defaultBranch == "" {
		var err error
		defaultBranch, _, err = defaultBranchFromRemote(ctx, storePath)
		if err != nil {
			return "", err
		}
		remoteChecked = true
	}
	if defaultBranch != "" {
		_, _ = gitcmd.Run(ctx, []string{"symbolic-ref", "refs/remotes/origin/HEAD", fmt.Sprintf("refs/remotes/origin/%s", defaultBranch)}, gitcmd.Options{Dir: storePath})
	}

	if fetch {
		if log {
			gitcmd.Logf("git fetch --prune")
		}
		if _, err := gitcmd.Run(ctx, []string{"fetch", "--prune"}, gitcmd.Options{Dir: storePath}); err != nil {
			return "", err
		}
	} else if remoteChecked {
		if err := touchFetchHead(storePath); err != nil {
			return "", err
		}
	}

	return defaultBranch, nil
}

func defaultBranchFromRemote(ctx context.Context, storePath string) (string, string, error) {
	res, err := gitcmd.Run(ctx, []string{"ls-remote", "--symref", "origin", "HEAD"}, gitcmd.Options{Dir: storePath})
	if err != nil {
		return "", "", err
	}
	branch, hash := coregitparse.ParseRemoteHeadSymref(res.Stdout)
	return branch, hash, nil
}

func localRemoteHash(ctx context.Context, storePath, branch string) (string, error) {
	ref := fmt.Sprintf("refs/remotes/origin/%s", branch)
	hash, exists, err := gitcmd.ShowRef(ctx, storePath, ref)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}
	return hash, nil
}

func localHeadHash(ctx context.Context, storePath, branch string) (string, error) {
	ref := fmt.Sprintf("refs/heads/%s", branch)
	hash, exists, err := gitcmd.ShowRef(ctx, storePath, ref)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}
	return hash, nil
}

func pruneLocalHeads(ctx context.Context, storePath, keepBranch string) error {
	worktreeBranches, _ := worktreeBranchNames(ctx, storePath)
	res, err := gitcmd.Run(ctx, []string{"show-ref", "--heads"}, gitcmd.Options{Dir: storePath})
	if err != nil && res.ExitCode != 1 {
		return err
	}
	lines := strings.Split(strings.TrimSpace(res.Stdout), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		ref := parts[1]
		if !strings.HasPrefix(ref, "refs/heads/") {
			continue
		}
		name := strings.TrimPrefix(ref, "refs/heads/")
		if name == keepBranch {
			continue
		}
		if _, ok := worktreeBranches[name]; ok {
			continue
		}
		_, _ = gitcmd.Run(ctx, []string{"update-ref", "-d", ref}, gitcmd.Options{Dir: storePath})
	}
	return nil
}

func worktreeBranchNames(ctx context.Context, storePath string) (map[string]struct{}, error) {
	out, err := gitcmd.WorktreeListPorcelain(ctx, storePath)
	if err != nil {
		return nil, err
	}
	branches := make(map[string]struct{})
	for _, name := range coregitparse.ParseWorktreeBranchNames(out) {
		branches[name] = struct{}{}
	}
	return branches, nil
}

func fetchGraceDuration() time.Duration {
	return corerepostore.FetchGraceDuration(os.Getenv("GION_FETCH_GRACE_SECONDS"))
}

func recentlyFetched(storePath string, grace time.Duration) bool {
	return corerepostore.RecentlyFetched(storePath, grace)
}

func touchFetchHead(storePath string) error {
	return corerepostore.TouchFetchHead(storePath)
}

func localDefaultBranch(ctx context.Context, storePath string) (string, error) {
	ref, ok, err := gitcmd.SymbolicRef(ctx, storePath, "refs/remotes/origin/HEAD")
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	if branch, ok := coregitref.ParseOriginHeadRef(ref); ok {
		return branch, nil
	}
	return "", nil
}
