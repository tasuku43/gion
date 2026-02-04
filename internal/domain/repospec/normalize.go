package repospec

import corerepospec "github.com/tasuku43/gion-core/repospec"

type Spec = corerepospec.Spec

func Normalize(input string) (Spec, error) {
	return corerepospec.Normalize(input)
}
