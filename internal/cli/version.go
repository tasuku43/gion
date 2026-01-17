package cli

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// These are intended to be set via -ldflags.
//
// Example:
//
//	go build -ldflags "-X github.com/tasuku43/gws/internal/cli.version=v0.1.0 -X github.com/tasuku43/gws/internal/cli.commit=abc123 -X github.com/tasuku43/gws/internal/cli.date=2026-01-17"
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func versionLine() string {
	v := strings.TrimSpace(version)
	if v == "" {
		v = "dev"
	}
	parts := []string{fmt.Sprintf("gws %s", v)}
	if c := strings.TrimSpace(commit); c != "" {
		parts = append(parts, c)
	}
	if d := strings.TrimSpace(date); d != "" {
		parts = append(parts, d)
	}
	parts = append(parts, fmt.Sprintf("(%s %s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	return strings.Join(parts, " ")
}

func printVersion(w io.Writer) {
	fmt.Fprintln(w, versionLine())
}
