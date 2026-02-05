package repospec

import corerepospec "github.com/tasuku43/gion-core/repospec"

type Spec = corerepospec.Spec

func Normalize(input string) (Spec, error) {
	return corerepospec.Normalize(input)
}

func DisplaySpec(input string) string {
	return corerepospec.DisplaySpec(input)
}

func DisplayName(input string) string {
	return corerepospec.DisplayName(input)
}

func SpecFromKey(repoKey string) string {
	return corerepospec.SpecFromKey(repoKey)
}
