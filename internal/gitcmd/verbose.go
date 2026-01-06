package gitcmd

import (
	"fmt"
	"os"
)

var verbose bool

func SetVerbose(v bool) {
	verbose = v
}

func IsVerbose() bool {
	return verbose
}

func Logf(format string, args ...any) {
	if verbose {
		return
	}
	fmt.Fprintf(os.Stderr, "\x1b[36m$ "+format+"\x1b[0m\n", args...)
}
