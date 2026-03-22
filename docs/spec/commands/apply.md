---
title: "gion apply"
status: implemented
---

## Synopsis
`gion apply [--root <path>] [--no-prompt]`

## Intent
Reconcile the filesystem to match `gion.yaml` by computing a diff, showing a plan, and applying the changes after confirmation.

## Behavior
- Loads `<root>/gion.yaml`; errors if missing or invalid.
- Scans `<root>/workspaces` to build the current state.
- Computes a plan with `add`, `remove`, and `update` actions:
  - `add`: workspace or repo entry exists in manifest but not on filesystem.
  - `remove`: exists on filesystem but not in manifest.
  - `update`: exists in both but differs by repo alias, repo key, or branch.
- Renders a human-readable plan summary before any changes (same format as `gion plan`).
- By default, prompts for confirmation if any changes exist.
  - `remove` actions are marked as destructive.
  - If only non-destructive adds are present, prompt can be skipped with `--no-prompt`.
  - For destructive actions, the prompt does not repeat per-repo git status output; users should review the plan output above before confirming.
- If confirmed, applies actions in a stable order: removes, then updates, then adds.
  - When a repo update is a branch rename only (same repo key, different branch), gion renames the branch in-place (no worktree remove/add) to match common local development workflows.
- During destructive workspace removal, if the current process cwd is inside `workspaces/<id>/...`, gion must first shift process cwd to `<root>` before any worktree removal starts.
  - When shell integration from `gion shell init` is active, successful workspace removal must emit a shell action that changes the parent shell cwd to `<root>`.
  - On failure, gion must not emit a shell action.
- During destructive workspace removal, gion must persist transient operation state at `<root>/.gion/state/operations/workspace-remove/<workspace-id>.json`.
  - The state file tracks which repo aliases were already removed and whether the workspace directory is still present.
  - The state file is operational metadata only; it is not part of the canonical manifest inventory and must not be imported into `gion.yaml`.
  - On successful workspace removal, gion must delete the state file.
- Before applying a destructive workspace removal, gion must check for an existing unfinished state file for that workspace.
  - If one exists, `gion apply` must fail fast before starting any new destructive work for that workspace.
  - The error must identify the workspace and summarize the partial progress so the user can inspect the unfinished removal state.
- When applying `add` actions that require creating a new branch:
  - If the target `branch` already exists in the bare store, gion checks it out when adding the worktree.
  - If the branch does not exist, gion creates it from:
    - `base_ref` if present in the repo entry in `gion.yaml`, otherwise
    - the repo's detected default branch (prefer `refs/remotes/origin/HEAD`), otherwise fallback heuristics (`HEAD`, then common branch names).
- When gion creates a new branch during apply, it records the chosen base as `base_branch` in the workspace `.gion/metadata.json` (workspace-level, optional) so a future `gion import` can restore `base_ref` in `gion.yaml`.
- Updates `gion.yaml` by rewriting the full file after successful apply.

## Output (IA)
- `Plan` section: plan summary (same as `gion plan`).
  - When interactive, the confirmation UI is rendered as a `Step` section between `Plan` and `Apply`.
  - Once confirmation completes, the `Step` prompt should not remain in the final visible transcript.
- `Apply` section: execution summary at the workspace / repo level.
  - Normal output should not dump raw git commands such as `git worktree add` / `git worktree remove`.
- `Result` section: completion summary (e.g. applied counts) and manifest rewrite note.

## Flags
- `--no-prompt`: skip confirmation (errors if any removals are present).

## Success Criteria
- Filesystem state matches the manifest.
- `gion.yaml` is rewritten to a normalized form.

## Failure Modes
- Manifest file missing or invalid.
- Filesystem or git errors while applying actions.
- `--no-prompt` used with destructive actions.
- Unfinished workspace removal state already exists at `<root>/.gion/state/operations/workspace-remove/<workspace-id>.json`.
