package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tasuku43/gion/internal/domain/manifest"
	"github.com/tasuku43/gion/internal/domain/preset"
)

func normalizeManifestAddArgs(args []string) []string {
	// Go's stdlib flag package stops parsing at the first non-flag argument.
	// For manifest add, we want to accept flags even when they appear after the
	// positional args, as the help text documents.
	//
	// Example:
	//   gion manifest add --repo <repo> <WORKSPACE_ID> --no-prompt
	// should behave the same as:
	//   gion manifest add --repo <repo> --no-prompt <WORKSPACE_ID>
	return normalizeArgsFlagsFirst(args, map[string]struct{}{
		"--preset":       {},
		"-preset":        {},
		"--repo":         {},
		"-repo":          {},
		"--branch":       {},
		"-branch":        {},
		"--base":         {},
		"-base":          {},
		"--workspace-id": {},
		"-workspace-id":  {},
	})
}

func normalizeManifestPresetAddArgs(args []string) []string {
	return normalizeArgsFlagsFirst(args, map[string]struct{}{
		"--repo": {},
		"-repo":  {},
	})
}

func normalizeArgsFlagsFirst(args []string, requiresValue map[string]struct{}) []string {
	if len(args) == 0 {
		return args
	}

	flagArgs := make([]string, 0, len(args))
	positionalArgs := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Respect the conventional end-of-flags marker. Everything after `--` is positional.
		if arg == "--" {
			positionalArgs = append(positionalArgs, arg)
			if i+1 < len(args) {
				positionalArgs = append(positionalArgs, args[i+1:]...)
			}
			break
		}

		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "=") {
				flagArgs = append(flagArgs, arg)
				continue
			}

			if _, ok := requiresValue[arg]; ok {
				if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
					flagArgs = append(flagArgs, arg+"=")
					continue
				}
				flagArgs = append(flagArgs, arg, args[i+1])
				i++
				continue
			}

			flagArgs = append(flagArgs, arg)
			continue
		}

		positionalArgs = append(positionalArgs, arg)
	}

	out := make([]string, 0, len(flagArgs)+len(positionalArgs))
	out = append(out, flagArgs...)
	out = append(out, positionalArgs...)
	return out
}

func loadPresetNames(rootDir string) ([]string, error) {
	file, err := preset.Load(rootDir)
	if err != nil {
		return nil, err
	}
	names := preset.Names(file)
	if len(names) == 0 {
		return nil, fmt.Errorf("no presets found in %s", filepath.Join(rootDir, manifest.FileName))
	}
	return names, nil
}
