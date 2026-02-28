---
title: completion
description: Generate shell completion script
status: implemented
---

# gion completion

Generate shell completion script for bash or zsh.

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
