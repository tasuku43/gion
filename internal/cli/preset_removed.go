package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func runPresetRemoved(args []string) error {
	fs := flag.NewFlagSet("preset", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var helpFlag bool
	fs.BoolVar(&helpFlag, "help", false, "show help")
	fs.BoolVar(&helpFlag, "h", false, "show help")
	fs.Usage = func() {
		printPresetHelp(os.Stdout)
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if helpFlag {
		printPresetHelp(os.Stdout)
		return nil
	}
	return fmt.Errorf("gwst preset is removed; use: gwst manifest preset")
}
