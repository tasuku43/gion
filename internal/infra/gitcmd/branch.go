package gitcmd

import (
	"context"
	"fmt"
	"strings"
)

func BranchMove(ctx context.Context, dir, fromBranch, toBranch string) error {
	from := strings.TrimSpace(fromBranch)
	to := strings.TrimSpace(toBranch)
	if from == "" {
		return fmt.Errorf("from branch is required")
	}
	if to == "" {
		return fmt.Errorf("to branch is required")
	}
	if from == to {
		return nil
	}
	_, err := Run(ctx, []string{"branch", "-m", from, to}, Options{Dir: dir, ShowOutput: true})
	return err
}
