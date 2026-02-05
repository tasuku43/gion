package cli

import (
	"fmt"
	"strings"
)

func normalizeGlobalArgs(args []string) ([]string, error) {
	if len(args) == 0 {
		return args, nil
	}

	boolFlags := map[string]struct{}{
		"--no-prompt": {},
		"--debug":     {},
		"--help":      {},
		"-h":          {},
		"--version":   {},
	}

	global := make([]string, 0, len(args))
	rest := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Respect the conventional end-of-flags marker. Everything after `--` is not a flag.
		if arg == "--" {
			rest = append(rest, arg)
			if i+1 < len(args) {
				rest = append(rest, args[i+1:]...)
			}
			break
		}

		// Global string flag: --root <path> or --root=<path>
		if arg == "--root" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag needs an argument: --root")
			}
			next := args[i+1]
			if next == "--" {
				return nil, fmt.Errorf("flag needs an argument: --root")
			}
			// Keep behavior strict: treat a following token that looks like a flag as missing.
			if strings.HasPrefix(next, "-") {
				return nil, fmt.Errorf("flag needs an argument: --root")
			}
			global = append(global, arg, next)
			i++
			continue
		}
		if strings.HasPrefix(arg, "--root=") {
			global = append(global, arg)
			continue
		}

		// Global bool flags (also allow --flag=true style).
		if _, ok := boolFlags[arg]; ok {
			global = append(global, arg)
			continue
		}
		if strings.HasPrefix(arg, "--no-prompt=") ||
			strings.HasPrefix(arg, "--debug=") ||
			strings.HasPrefix(arg, "--help=") ||
			strings.HasPrefix(arg, "--version=") {
			global = append(global, arg)
			continue
		}

		rest = append(rest, arg)
	}

	out := make([]string, 0, len(global)+len(rest))
	out = append(out, global...)
	out = append(out, rest...)
	return out, nil
}
