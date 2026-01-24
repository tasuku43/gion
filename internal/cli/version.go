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
//	go build -ldflags "-X <module>/internal/cli.version=v0.1.0 -X <module>/internal/cli.commit=abc123 -X <module>/internal/cli.date=2026-01-17"
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func versionLine() string {
	return versionLineFor("gion")
}

func printVersion(w io.Writer) {
	fmt.Fprintln(w, versionLine())
}

func versionLineFor(name string) string {
	v := strings.TrimSpace(version)
	if v == "" {
		v = "dev"
	}
	program := strings.TrimSpace(name)
	if program == "" {
		program = "gion"
	}
	parts := []string{fmt.Sprintf("%s %s", program, v)}
	if c := strings.TrimSpace(commit); c != "" {
		parts = append(parts, c)
	}
	if d := strings.TrimSpace(date); d != "" {
		parts = append(parts, d)
	}
	parts = append(parts, fmt.Sprintf("(%s %s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	return strings.Join(parts, " ")
}

func printVersionFor(w io.Writer, name string) {
	fmt.Fprintln(w, versionLineFor(name))
}
