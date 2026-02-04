package repo

import (
	"path/filepath"
	"strings"

	"github.com/tasuku43/gion/internal/domain/repospec"
	"github.com/tasuku43/gion/internal/infra/paths"
)

// Spec is the normalized repo specification.
type Spec = repospec.Spec

// StorePath returns the path to the bare repo store for the spec.
func StorePath(rootDir string, spec repospec.Spec) string {
	return filepath.Join(paths.BareRoot(rootDir), spec.Host, spec.Owner, spec.Repo+".git")
}

// Normalize trims and validates a repo spec, returning the spec and trimmed input.
func Normalize(input string) (repospec.Spec, string, error) {
	trimmed := strings.TrimSpace(input)
	spec, err := repospec.Normalize(trimmed)
	if err != nil {
		return repospec.Spec{}, "", err
	}
	return spec, trimmed, nil
}

// DisplaySpec returns a normalized display string for a repo spec.
func DisplaySpec(input string) string {
	return repospec.DisplaySpec(input)
}

// DisplayName returns the repo name for display.
func DisplayName(input string) string {
	return repospec.DisplayName(input)
}

// SpecFromKey converts a repo key (host/owner/repo.git) into a cloneable spec.
func SpecFromKey(repoKey string) string {
	return repospec.SpecFromKey(repoKey)
}
