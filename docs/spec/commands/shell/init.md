---
title: "gion shell init"
status: implemented
---

## Synopsis
`gion shell init [shell] [--with-completion[=true|false]]`

## Intent
Print shell integration code that allows `gion` to apply parent-shell side effects through an action-file protocol.

## Behavior
- Supports `bash` and `zsh`.
- If `shell` is omitted, infer it from `SHELL`; default to `zsh` when inference fails.
- Prints a shell wrapper function to stdout.
- The wrapper:
  - creates a temporary action file
  - runs `command gion` with `GION_SHELL_ACTION_FILE` pointing at that file
  - applies the action file only when the command succeeds
- `--with-completion` appends the same shell-specific script as `gion shell completion <shell>`.
- `--with-completion=<value>` accepts `true/false`, `1/0`, `yes/no`, `on/off`.

## Examples
- `eval "$(gion shell init zsh)"`
- `eval "$(gion shell init zsh --with-completion)"`

## Failure Modes
- Unsupported shell.
- Invalid `--with-completion` value.
