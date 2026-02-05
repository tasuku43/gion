package repo

import (
	corerepostore "github.com/tasuku43/gion-core/repostore"
	"github.com/tasuku43/gion/internal/infra/paths"
)

type Entry struct {
	RepoKey   string
	StorePath string
}

func List(rootDir string) ([]Entry, []error, error) {
	entries, warnings, err := corerepostore.List(paths.BareRoot(rootDir))
	if err != nil {
		return nil, warnings, err
	}
	result := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, Entry{
			RepoKey:   entry.RepoKey,
			StorePath: entry.StorePath,
		})
	}
	return result, warnings, nil
}
