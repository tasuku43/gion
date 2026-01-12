package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/tasuku43/gws/internal/core/output"
	"github.com/tasuku43/gws/internal/domain/workspace"
	"github.com/tasuku43/gws/internal/ui"
)

func runWorkspaceOpen(ctx context.Context, rootDir string, args []string, noPrompt bool) error {
	openFlags := flag.NewFlagSet("open", flag.ContinueOnError)
	var helpFlag bool
	var shellFlag bool
	openFlags.BoolVar(&shellFlag, "shell", false, "spawn interactive shell")
	openFlags.BoolVar(&helpFlag, "help", false, "show help")
	openFlags.BoolVar(&helpFlag, "h", false, "show help")
	openFlags.SetOutput(os.Stdout)
	openFlags.Usage = func() {
		printOpenHelp(os.Stdout)
	}
	if err := openFlags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if helpFlag {
		printOpenHelp(os.Stdout)
		return nil
	}
	if openFlags.NArg() > 1 {
		return fmt.Errorf("usage: gws open [<WORKSPACE_ID>] [--shell]")
	}

	workspaceID := ""
	if openFlags.NArg() == 1 {
		workspaceID = openFlags.Arg(0)
	}

	if workspaceID == "" {
		if noPrompt {
			return fmt.Errorf("workspace id is required without prompt")
		}
		workspaces, wsWarn, err := workspace.List(rootDir)
		if err != nil {
			return err
		}
		if len(wsWarn) > 0 {
			// ignore warnings for selection
		}
		workspaceChoices := buildWorkspaceChoices(ctx, workspaces)
		if len(workspaceChoices) == 0 {
			return fmt.Errorf("no workspaces found")
		}
		theme := ui.DefaultTheme()
		useColor := isatty.IsTerminal(os.Stdout.Fd())
		workspaceID, err = ui.PromptWorkspace("gws open", workspaceChoices, theme, useColor)
		if err != nil {
			return err
		}
	}

	wsDir := filepath.Join(rootDir, "workspaces", workspaceID)
	if info, err := os.Stat(wsDir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("workspace does not exist: %s", wsDir)
		}
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("workspace path is not a directory: %s", wsDir)
	}

	shellPath := strings.TrimSpace(os.Getenv("SHELL"))
	cmdPath, cmdArgs := shellCommandForOpen(shellPath)
	cmdDisplay := cmdPath
	if len(cmdArgs) > 0 {
		cmdDisplay = fmt.Sprintf("%s %s", cmdPath, strings.Join(cmdArgs, " "))
	}
	shellInfo := cmdPath
	if len(cmdArgs) > 0 {
		shellInfo = fmt.Sprintf("%s (interactive)", cmdPath)
	}
	theme := ui.DefaultTheme()
	useColor := isatty.IsTerminal(os.Stdout.Fd())
	renderer := ui.NewRenderer(os.Stdout, theme, useColor)
	output.SetStepLogger(renderer)
	defer output.SetStepLogger(nil)

	renderer.Section("Info")
	renderer.Bullet(fmt.Sprintf("workspace: %s", workspaceID))
	renderer.Bullet(fmt.Sprintf("path: %s", wsDir))
	renderer.Bullet(fmt.Sprintf("shell: %s", shellInfo))
	renderer.Bullet("note: subshell; parent cwd unchanged")
	renderer.Blank()
	startSteps(renderer)
	output.Step("chdir")
	output.Log(wsDir)
	output.Step("launch subshell")
	output.Log(cmdDisplay)
	renderer.Blank()
	renderer.Section("Result")
	renderer.Bullet("enter subshell (type `exit` to return)")
	if err := os.Chdir(wsDir); err != nil {
		return fmt.Errorf("chdir workspace: %w", err)
	}
	cmd := exec.CommandContext(ctx, cmdPath, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("GWS_WORKSPACE=%s", workspaceID))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("open shell: %w", err)
	}
	return nil
}

func shellCommandForOpen(shellPath string) (string, []string) {
	if strings.TrimSpace(shellPath) == "" {
		shellPath = "/bin/sh"
	}
	name := filepath.Base(shellPath)
	if isInteractiveShell(name) {
		return shellPath, []string{"-i"}
	}
	return shellPath, nil
}

func isInteractiveShell(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "bash", "zsh", "sh", "fish", "ksh", "dash", "tcsh", "csh":
		return true
	default:
		return false
	}
}
