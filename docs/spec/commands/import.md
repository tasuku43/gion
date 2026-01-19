---
title: "gwst import"
status: planned
---

## Synopsis
`gwst import [--root <path>] [--no-prompt]`

## Intent
Rebuild `manifest.yaml` from the filesystem and `.gwst/metadata.json` to restore the canonical workspace inventory.

## Behavior
- Scans `<root>/workspaces` to build the current filesystem state.
- For each workspace:
  - Loads `.gwst/metadata.json` when present to restore `mode`, `description`, `template_name`, `source_url`.
  - Derives repo branches from each worktree's Git state.
- Rewrites `<root>/manifest.yaml` as a whole, reflecting the current filesystem state.
- By default, shows a summary of changes and prompts for confirmation.
  - `--no-prompt` skips confirmation.

## Output
- Prints the summary of imported workspaces and repos.
- Reports warnings for unreadable workspaces or invalid metadata.

## Success Criteria
- `manifest.yaml` reflects the current filesystem state.

## Failure Modes
- Root directory missing or inaccessible.
- Filesystem errors while scanning workspaces.
- Invalid metadata that prevents import (reported as warnings; fatal only if no valid workspaces remain).
