---
title: "gws ls"
status: implemented
---

## Synopsis
`gws ls`

## Intent
List workspaces under `<root>/workspaces` and show a quick view of the repos attached to each.

## Behavior
- Scans `<root>/workspaces` for directories; ignores non-directories.
- For each workspace, scans its contents to discover repo worktrees (alias, repo key, branch, path) and renders them in a tree view.
- Collects and reports non-fatal warnings from scanning workspaces or repos.

## Success Criteria
- Existing workspaces are listed; command succeeds even if none exist (empty result).

## Failure Modes
- Root path inaccessible or `workspaces/` is not a directory.
- Filesystem or git errors while scanning workspaces (reported as warnings; unrecoverable errors fail the command).
