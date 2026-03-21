package shellaction

import (
	"fmt"
	"os"
	"strings"
)

const FileEnv = "GION_SHELL_ACTION_FILE"

func EmitCD(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	target := strings.TrimSpace(os.Getenv(FileEnv))
	if target == "" {
		return nil
	}
	line := fmt.Sprintf("builtin cd -- %s\n", shellQuote(path))
	if err := os.WriteFile(target, []byte(line), 0o600); err != nil {
		return fmt.Errorf("write shell action file: %w", err)
	}
	return nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}
