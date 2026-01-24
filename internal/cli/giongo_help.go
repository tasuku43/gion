package cli

import (
	"fmt"
	"io"
)

func printGiongoHelp(w io.Writer) {
	fmt.Fprintln(w, `giongo - interactive workspace/worktree picker

Usage:
  giongo [--print] [--root <path>]

Options:
  --print         print the selected absolute path
  --root <path>   override root directory
  -h, --help      show help
  --version       print version`)
}
