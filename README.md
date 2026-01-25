# gion — Git workspaces (built on Git worktrees) as code, with guardrails.

Git workspaces as code, with guardrails.  
Built on Git worktrees: define a YAML inventory, then plan/apply to reconcile safely.

## Demo (45s)

Coming soon (new demo video in progress).

This demo shows:
- Create many workspaces at once
- Jump fast with `giongo`
- Clean up in bulk with guardrails (plan warns before deletions)

## Quickstart (5 minutes)

### 1) Install

```bash
brew tap tasuku43/gion
brew install gion
```

Other options:
- Version pinning with mise (optional): `mise use -g github:tasuku43/gion@<version>`
- Manual install via GitHub Releases (download archive → put `gion` on your PATH)
- Build from source requires Go 1.24+

### 2) Prepare a repo store

```bash
gion repo get git@github.com:org/backend.git
```

### 3) Create a workspace (includes Plan + Apply by default)

```bash
gion manifest add --repo git@github.com:org/backend.git PROJ-123
```

### 4) Jump into a workspace (interactive)

```bash
cd "$(giongo --print)"
```

### 5) Remove safely (guardrails on by default)

```bash
gion manifest rm PROJ-123
```

## What gion solves

- **Reproducibility:** workspace inventory is code (`gion.yaml`), synced safely via diffs
- **Scale:** create/remove many worktrees without losing track
- **Safety:** `plan` warns about risky deletions (dirty/unpushed/diverged/unknown)
- **Bulk operations without fear:** visible guardrails for mass create/delete
- **Multi-repo tasks:** group multiple repos under one workspace (see “From presets” below)

## How it works (Declare → Diff → Reconcile)

- **Declare** desired workspaces in `gion.yaml` (inventory / desired state)
- **Diff** with `gion plan` (read-only)
- **Reconcile** with `gion apply` (shows a plan; confirms destructive changes)

## Create in bulk: entry points

### From PRs / issues (GitHub only)

Interactive bulk selection (multi-select in the picker):

```bash
gion manifest add
```

Notes:
- Requires `gh` (authenticated) to fetch metadata.
- The picker supports bulk selection of PRs/issues, then a single apply.

Direct URL (single workspace):

```bash
gion manifest add --review https://github.com/owner/repo/pull/123
gion manifest add --issue  https://github.com/owner/repo/issues/123
```

### From presets (multi-repo “task workspace”)

```yaml
presets:
  app:
    repos:
      - git@github.com:org/backend.git
      - git@github.com:org/frontend.git
      - git@github.com:org/infra.git
```

```bash
gion manifest add --preset app PROJ-123
```

## Details (create → move → delete)

### Create

Interactive front-end to the inventory:

```bash
gion manifest add
```

Run with flags to skip prompts:

```bash
gion manifest add --repo git@github.com:org/backend.git PROJ-123
```

### Move fast with giongo

`giongo` is a small companion binary that jumps into a workspace or repo using a picker.  
It does not change any state.

Example (zsh function):

```bash
giongo() {
  if [[ "$1" == "init" || "$1" == "--help" || "$1" == "-h" || "$1" == "--version" || "$1" == "--print" ]]; then
    command giongo "$@"
    return $?
  fi
  local dest
  dest="$(command giongo --print "$@")" || return $?
  [[ -n "$dest" ]] && cd "$dest"
}
```

Shortcut (auto-generate the function for your shell):

```bash
eval "$(giongo init)"
```

Notes:
- `giongo init` outputs a bash/zsh function definition.
- For a permanent setup, paste the output into `~/.zshrc` or `~/.bashrc`.

### Edit the inventory directly (AI-friendly)

`gion.yaml` is just YAML. You can edit it directly (humans or AI), then review/apply changes:

```bash
gion plan
gion apply
```

Minimal example:

```yaml
version: 1
workspaces:
  PROJ-123:
    description: "fix login flow"
    mode: repo
    repos:
      - alias: backend
        repo_key: github.com/org/backend.git
        branch: PROJ-123
```

Notes:
- `gion.yaml` is gion-managed: commands rewrite the whole file, so comments/ordering may not be preserved.
- If your local manual changes are the source of truth, use `gion import` to follow them back into `gion.yaml`.

### Cleanup

Manual removal (explicit human judgment):

```bash
gion manifest rm PROJ-123
```

Automatic cleanup (conservative):

```bash
gion manifest gc
```

`gion manifest gc` removes workspace entries from `gion.yaml` only when they are highly likely safe to delete, then (by default) runs `gion apply` to reconcile.

### Import

If the filesystem is the source of truth, rebuild the inventory:

```bash
gion import
```

## Requirements

- Git
- `gh` CLI (optional; required for `gion manifest add --review` and `gion manifest add --issue` — GitHub only)

## Help and docs

- `docs/README.md` for documentation index
- `docs/spec/README.md` for specs index and status
- `docs/spec/commands/` for per-command specs
- `docs/spec/core/INVENTORY.md` for `gion.yaml` format
- `docs/spec/core/PRESETS.md` for preset format
- `docs/spec/core/DIRECTORY_LAYOUT.md` for the file layout
- `docs/spec/ui/UI.md` for output conventions
- `docs/concepts/CONCEPT.md` for background and motivation

## Maintainer

- @tasuku43
