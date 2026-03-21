package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runShell(args []string) error {
	if len(args) == 0 {
		printShellHelp(os.Stdout)
		return nil
	}

	switch strings.TrimSpace(args[0]) {
	case "-h", "--help", "help":
		printShellHelp(os.Stdout)
		return nil
	case "init":
		return runShellInit(args[1:])
	case "completion":
		return runShellCompletion(args[1:])
	default:
		return fmt.Errorf("unknown shell command: %s", args[0])
	}
}

func runShellInit(args []string) error {
	shellName := ""
	withCompletion := false
	rest := append([]string{}, args...)
	for len(rest) > 0 {
		cur := strings.TrimSpace(rest[0])
		switch cur {
		case "-h", "--help", "help":
			printShellHelp(os.Stdout)
			return nil
		case "--with-completion":
			withCompletion = true
			rest = rest[1:]
		default:
			if strings.HasPrefix(cur, "--with-completion=") {
				value := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(cur, "--with-completion=")))
				switch value {
				case "1", "true", "yes", "on":
					withCompletion = true
				case "0", "false", "no", "off":
					withCompletion = false
				default:
					return fmt.Errorf("invalid --with-completion value: %q (supported: true/false)", value)
				}
				rest = rest[1:]
				continue
			}
			if strings.HasPrefix(cur, "-") {
				return fmt.Errorf("unknown flag for shell init: %q", cur)
			}
			if shellName != "" {
				return fmt.Errorf("unexpected args for shell init: %q", strings.Join(rest, " "))
			}
			shellName = cur
			rest = rest[1:]
		}
	}

	if shellName == "" {
		shellName = detectShellName()
	}
	if shellName == "" {
		shellName = "zsh"
	}

	script, err := renderShellInitScript(shellName, withCompletion)
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, script)
	return nil
}

func runShellCompletion(args []string) error {
	shellName := ""
	if len(args) > 1 {
		return fmt.Errorf("unexpected args for shell completion: %q", strings.Join(args[1:], " "))
	}
	if len(args) == 1 {
		shellName = strings.TrimSpace(args[0])
	} else {
		shellName = detectShellName()
	}
	if shellName == "" {
		shellName = "zsh"
	}
	return runCompletion([]string{shellName})
}

func detectShellName() string {
	raw := strings.TrimSpace(os.Getenv("SHELL"))
	if raw == "" {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(filepath.Base(raw)))
}

func renderShellInitScript(shellName string, withCompletion bool) (string, error) {
	initScript, err := renderPOSIXShellInitScript(shellName)
	if err != nil {
		return "", err
	}
	if !withCompletion {
		return initScript, nil
	}
	completionScript, err := renderShellCompletionScript(shellName)
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(initScript, "\n") {
		return initScript + "\n" + completionScript, nil
	}
	return initScript + "\n\n" + completionScript, nil
}

func renderPOSIXShellInitScript(shellName string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(shellName)) {
	case "zsh", "bash", "sh":
	default:
		return "", fmt.Errorf("unsupported shell: %s (supported: %s)", shellName, SupportedShells)
	}
	return fmt.Sprintf(`# gion shell integration (%s)
# Add this line to your shell rc file, then restart the shell:
#   eval "$(gion shell init %s)"
gion() {
  local __gion_action_file __gion_status
  __gion_action_file="$(mktemp "${TMPDIR:-/tmp}/gion-shell-action.XXXXXX")" || return 1
  GION_SHELL_ACTION_FILE="$__gion_action_file" command gion "$@"
  __gion_status=$?
  if [ $__gion_status -ne 0 ]; then
    rm -f "$__gion_action_file"
    return $__gion_status
  fi
  if [ -s "$__gion_action_file" ]; then
    eval "$(cat "$__gion_action_file")"
  fi
  rm -f "$__gion_action_file"
}
`, shellName, shellName), nil
}
