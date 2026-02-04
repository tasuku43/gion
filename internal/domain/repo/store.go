package repo

import (
	"context"
	"fmt"
	"os"

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
	result, err := corerepostore.EnsureStore(ctx, storeAccessAdapter{}, corerepostore.EnsureStoreRequest{
		RepoKey:       spec.RepoKey,
		RemoteURL:     remoteURL,
		StorePath:     storePath,
		RepoSpec:      repo,
		MustExist:     false,
		Fetch:         false,
		FetchGraceEnv: os.Getenv("GION_FETCH_GRACE_SECONDS"),
		Log:           true,
	})
	if err != nil {
		return Store{}, err
	}
	return fromCoreStore(result.Store), nil
}

func Open(ctx context.Context, rootDir string, repo string, fetch bool) (Store, error) {
	spec, remoteURL, err := Normalize(repo)
	if err != nil {
		return Store{}, err
	}

	storePath := storePathForSpec(rootDir, spec)
	result, err := corerepostore.EnsureStore(ctx, storeAccessAdapter{}, corerepostore.EnsureStoreRequest{
		RepoKey:       spec.RepoKey,
		RemoteURL:     remoteURL,
		StorePath:     storePath,
		RepoSpec:      repo,
		MustExist:     true,
		Fetch:         fetch,
		FetchGraceEnv: os.Getenv("GION_FETCH_GRACE_SECONDS"),
		Log:           true,
	})
	if err != nil {
		return Store{}, err
	}
	return fromCoreStore(result.Store), nil
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
	_, err := corerepostore.NormalizeStore(ctx, normalizerGitAdapter{}, storePath, fetch, os.Getenv("GION_FETCH_GRACE_SECONDS"), true)
	return err
}

type storeAccessAdapter struct{}

func (storeAccessAdapter) DirExists(path string) (bool, error) {
	return paths.DirExists(path)
}

func (storeAccessAdapter) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (storeAccessAdapter) CloneBare(ctx context.Context, remoteURL, storePath string) error {
	gitcmd.Logf("git clone --bare %s %s", remoteURL, storePath)
	_, err := gitcmd.Run(ctx, []string{"clone", "--bare", remoteURL, storePath}, gitcmd.Options{})
	return err
}

func (storeAccessAdapter) NormalizeStore(ctx context.Context, storePath string, fetch bool, fetchGraceEnv string, log bool) error {
	_, err := corerepostore.NormalizeStore(ctx, normalizerGitAdapter{}, storePath, fetch, fetchGraceEnv, log)
	return err
}

func fromCoreStore(store corerepostore.Store) Store {
	return Store{
		RepoKey:   store.RepoKey,
		StorePath: store.StorePath,
		RemoteURL: store.RemoteURL,
	}
}

func ensureDefaultBranch(ctx context.Context, storePath string, fetch bool, log bool) (string, error) {
	return corerepostore.EnsureDefaultBranch(ctx, normalizerGitAdapter{}, storePath, fetch, os.Getenv("GION_FETCH_GRACE_SECONDS"), log)
}

type normalizerGitAdapter struct{}

func (normalizerGitAdapter) ConfigureRemoteFetch(ctx context.Context, storePath string) error {
	if _, err := gitcmd.Run(ctx, []string{"config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*"}, gitcmd.Options{Dir: storePath}); err != nil {
		return err
	}
	return nil
}

func (normalizerGitAdapter) LocalDefaultBranch(ctx context.Context, storePath string) (string, error) {
	return localDefaultBranch(ctx, storePath)
}

func (normalizerGitAdapter) DefaultBranchFromRemote(ctx context.Context, storePath string) (string, error) {
	branch, _, err := defaultBranchFromRemote(ctx, storePath)
	return branch, err
}

func (normalizerGitAdapter) SetRemoteHead(ctx context.Context, storePath, branch string) error {
	if _, err := gitcmd.Run(ctx, []string{"symbolic-ref", "refs/remotes/origin/HEAD", fmt.Sprintf("refs/remotes/origin/%s", branch)}, gitcmd.Options{Dir: storePath}); err != nil {
		return err
	}
	return nil
}

func (normalizerGitAdapter) FetchPrune(ctx context.Context, storePath string, log bool) error {
	if log {
		gitcmd.Logf("git fetch --prune")
	}
	if _, err := gitcmd.Run(ctx, []string{"fetch", "--prune"}, gitcmd.Options{Dir: storePath}); err != nil {
		return err
	}
	return nil
}

func (normalizerGitAdapter) WorktreeBranches(ctx context.Context, storePath string) ([]string, error) {
	out, err := gitcmd.WorktreeListPorcelain(ctx, storePath)
	if err != nil {
		return nil, err
	}
	return coregitparse.ParseWorktreeBranchNames(out), nil
}

func (normalizerGitAdapter) HeadRefs(ctx context.Context, storePath string) ([]string, error) {
	res, err := gitcmd.Run(ctx, []string{"show-ref", "--heads"}, gitcmd.Options{Dir: storePath})
	if err != nil && res.ExitCode != 1 {
		return nil, err
	}
	return coregitparse.ParseHeadRefs(res.Stdout), nil
}

func (normalizerGitAdapter) DeleteRef(ctx context.Context, storePath, ref string) error {
	_, err := gitcmd.Run(ctx, []string{"update-ref", "-d", ref}, gitcmd.Options{Dir: storePath})
	return err
}

func (normalizerGitAdapter) TouchFetchHead(storePath string) error {
	return corerepostore.TouchFetchHead(storePath)
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
