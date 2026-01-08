---
title: "gws rm"
status: implemented
---

## Synopsis
`gws rm [<WORKSPACE_ID>]`

## Intent
Safely remove a workspace and all of its worktrees, refusing when repositories are dirty or status cannot be read.

## Behavior
- With `WORKSPACE_ID` provided: targets that workspace.
- Without it: scans workspaces, classifies each as removable or blocked (dirty or status errors), and prompts the user to choose one of the removable entries. Fails if none are removable.
- Before removal, gathers warnings (e.g., ahead-of-upstream, missing upstream, status errors) and displays them.
- Calls `workspace.Remove`, which:
  - Validates the workspace exists.
  - Fails if any repo has uncommitted/untracked/unstaged/unmerged changes.
  - Runs `git worktree remove <worktree>` for each repo’s worktree.
  - Deletes the workspace directory.

## Success Criteria
- Workspace directory no longer exists; associated worktrees are removed from their bare stores.

## Failure Modes
- Workspace not found.
- Dirty worktrees or status errors block removal.
- Git errors while removing worktrees.
- Filesystem errors while deleting the workspace directory.
