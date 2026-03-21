---
title: completion
description: Generate shell completion script
status: implemented
---

# gion completion

Generate shell completion script for bash or zsh.

The emitted script also installs a lightweight shell wrapper around `gion` so successful commands can apply shell actions such as `cd` after destructive workspace removal.

## Usage

```
gion completion [shell]
```

If no shell is specified, defaults to `bash`.

### Shells

| Shell | Description |
|-------|-------------|
| bash  | Generate bash completion |
| zsh   | Generate zsh completion |

## Setup

### bash

Add to `~/.bashrc`:

```bash
eval "$(gion completion bash)"
```

### zsh

Add to `~/.zshrc`:

```zsh
eval "$(gion completion zsh)"
```

## Shell action behavior

When the completion script wrapper is active:

- `gion` runs through a shell function wrapper rather than calling the binary directly
- the wrapper provides a temporary action-file path via `GION_SHELL_ACTION_FILE`
- on success, if `gion` writes a shell action, the wrapper sources it in the parent shell

Current use:

- successful workspace removal can move the parent shell cwd back to `GION_ROOT` when the previous cwd was inside the removed workspace

## Completion Coverage

### Commands

- `init`, `doctor`, `repo`, `manifest`, `plan`, `import`, `apply`, `version`, `help`, `completion`
- Aliases: `man`, `m` (for manifest)

### Subcommands

| Command | Subcommands |
|---------|-------------|
| repo    | `get`, `ls`, `rm` |
| manifest | `ls`, `add`, `rm`, `gc`, `validate`, `preset` |
| manifest preset | `ls`, `add`, `rm`, `validate` |

### Flags

Global flags:
- `--root`, `--no-prompt`, `--debug`, `--help`, `--version`

Command-specific flags are also completed (e.g., `--fix`, `--self` for doctor).
