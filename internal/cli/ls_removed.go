package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func runLsRemoved(args []string) error {
	fs := flag.NewFlagSet("ls", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var helpFlag bool
	fs.BoolVar(&helpFlag, "help", false, "show help")
	fs.BoolVar(&helpFlag, "h", false, "show help")
	_ = fs.Bool("details", false, "show git status details (removed)")
	fs.Usage = func() {
		printLsHelp(os.Stdout)
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if helpFlag {
		printLsHelp(os.Stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("usage: gwst ls [--details]")
	}
	return fmt.Errorf("gwst ls is removed; use: gwst manifest ls")
}
