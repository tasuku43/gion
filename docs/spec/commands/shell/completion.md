---
title: "gion shell completion"
status: implemented
---

## Synopsis
`gion shell completion [shell]`

## Intent
Print shell completion script for `gion`.

## Behavior
- Supports `bash` and `zsh`.
- If `shell` is omitted, infer it from `SHELL`; default to `zsh` when inference fails.
- Prints completion script to stdout.
- Does not include shell integration wrapper logic; use `gion shell init` for parent-shell side effects.

## Examples
- `source <(gion shell completion bash)`
- `source <(gion shell completion zsh)`

## Failure Modes
- Unsupported shell.
