# gws - Git Workspaces for Human + Agentic Development

gws moves local development from "clone-directory centric" to "workspace centric"
so humans and multiple AI agents can work in parallel without stepping on each other.

## Why gws

- In the era of AI agents, multiple actors edit in parallel and context collisions become common.
- gws promotes directories into explicit workspaces and manages them safely with Git worktrees.
- It focuses on creating, listing, and safely cleaning up work environments.

## What makes gws different

### 1) `create` is the center

One command, four creation modes:

```bash
gws create --repo git@github.com:org/repo.git
gws create --template app PROJ-123
gws create --review https://github.com/owner/repo/pull/123   # GitHub only
gws create --issue https://github.com/owner/repo/issues/123  # GitHub only
```

### 2) Template = pseudo-monorepo workspace

Define multiple repos as one task unit, then create them together:

```yaml
templates:
  app:
    repos:
      - git@github.com:org/api.git
      - git@github.com:org/web.git
```

```bash
gws create --template app PROJ-123
```

### 3) Guardrails on cleanup

`gws rm` refuses or asks for confirmation when workspaces are dirty, unpushed, or unknown:

```bash
gws rm PROJ-123
```

## Requirements

- Git
- Go 1.24+ (build/run from source)
- gh CLI (optional; required for `gws create --review` and `gws create --issue` — GitHub only)

## Install

Recommended:

```bash
brew tap tasuku43/gws
brew install gws
```

Version pinning (recommended):

```bash
mise use -g github:tasuku43/gws@v0.1.0
```

For details and other options, see `docs/guides/INSTALL.md`.

## Quickstart (5 minutes)

### 1) Initialize the root

```bash
gws init
```

This creates `GWS_ROOT` with the standard layout and a starter `templates.yaml`.

### 2) Define templates

Edit `templates.yaml` and list the repos you want in a workspace:

```yaml
templates:
  example:
    repos:
      - git@github.com:octocat/Hello-World.git
      - git@github.com:octocat/Spoon-Knife.git
```

Validate the file:

```bash
gws template validate
```

### 3) Fetch repos (bare store)

```bash
gws repo get git@github.com:octocat/Hello-World.git
gws repo get git@github.com:octocat/Spoon-Knife.git
```

### 4) Create a workspace

```bash
gws create --template example MY-123
```

Or create from a single repo:

```bash
gws create --repo git@github.com:octocat/Hello-World.git
```

Or run `gws create` with no args to pick a mode and fill inputs interactively.

### 5) Work and clean up

```bash
gws ls
gws open MY-123
gws status MY-123
gws rm MY-123
```

gws opens an interactive subshell at the workspace root.

## Provider support (summary)
- `gws create --repo` and `gws create --template` are provider-agnostic (any Git host URL).
- `gws create --review` and `gws create --issue` are GitHub-only today.

## How gws lays out files

gws keeps two top-level directories under `GWS_ROOT`:

```
GWS_ROOT/
  bare/        # bare repo store (shared Git objects)
  workspaces/  # task worktrees (one folder per workspace id)
  templates.yaml
```

Notes:

- Workspace id must be a valid Git branch name, and it becomes the worktree branch name.
- gws never changes your shell directory automatically.

## Root resolution

gws resolves `GWS_ROOT` in this order:

1. `--root <path>`
2. `GWS_ROOT` environment variable
3. `~/gws`

## Command overview (short)

- `gws init` - initialize root and `templates.yaml`
- `gws repo get <repo>` - fetch bare repo store
- `gws create ...` - create workspaces (repo/template/review/issue)
- `gws open [<id>]` - open a workspace in a subshell
- `gws status [<id>]` - check status
- `gws rm [<id>]` - remove workspace with guardrails

## Help and docs

- `docs/README.md` for documentation index
- `docs/spec/README.md` for specs index and status
- `docs/spec/core/TEMPLATES.md` for template format
- `docs/spec/core/DIRECTORY_LAYOUT.md` for the file layout
- `docs/spec/ui/UI.md` for output conventions
- `docs/concepts/CONCEPT.md` for the background and motivation

## Maintainer

- @tasuku43
