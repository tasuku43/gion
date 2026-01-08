---
title: "gws doctor"
status: implemented
---

## Synopsis
`gws doctor [--fix]`

## Intent
Detect common problems that block gws from working and surface them before users run other commands.

## Behavior
- Validates that a root directory was resolved.
- Checks the root layout for the presence of `bare/`, `src/`, `workspaces/`, and `templates.yaml`, reporting missing or invalid entries as issues.
- Scans existing workspaces and aggregates any warnings emitted while inspecting their repositories (e.g., unreadable worktrees).
- Lists repo stores and flags any store whose `origin` remote is missing or lacks a URL (`missing_remote`).
- `--fix` currently performs the same checks and returns the list of issues; no automatic fixes are applied yet (the `fixed` list remains empty).

## Success Criteria
- Command completes without errors; issues/warnings are printed for user action.

## Failure Modes
- Root directory not provided or inaccessible.
- Filesystem or git errors while inspecting workspaces or repo stores.
