package cli

import "github.com/tasuku43/gwst/internal/ui"

func buildPresetRepoChoices(rootDir string) ([]ui.PromptChoice, error) {
	return buildManifestPresetRepoChoices(rootDir)
}
