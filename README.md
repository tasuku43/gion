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

If you omit options, gws switches to an interactive flow:

```
$ gws create
Inputs
  • mode: s (type to filter)
    └─ repo - 1 repo only
    └─ issue - From an issue (multi-select, GitHub only)
    └─ review - From a review request (multi-select, GitHub only)
    └─ template - From template
```

Review/issue modes are also interactive (repo + multi-select):

```
$ gws create --review
Inputs
  • repo: org/gws
  • pull request: s (type to filter)
Info
  • selected
    └─ #123 Fix status output
    └─ #120 Add repo prompt
```

```
$ gws create --issue
Inputs
  • repo: org/gws
  • issue: s (type to filter)
Info
  • selected
    └─ #45 Improve template flow
    └─ #39 Add doctor checks
```

### 2) Template = pseudo-monorepo workspace

Define multiple repos as one task unit, then create them together:

```yaml
templates:
  app:
    repos:
      - git@github.com:org/backend.git
      - git@github.com:org/frontend.git
      - git@github.com:org/manifests.git
      - git@github.com:org/docs.git
```

```bash
gws create --template app PROJ-123
```

### 3) Guardrails on cleanup

`gws rm` refuses or asks for confirmation when workspaces are dirty, unpushed, or unknown:

```bash
gws rm PROJ-123
```

Omitting the workspace id prompts selection:

```
$ gws rm
Inputs
  • workspace: s (type to filter)
    └─ PROJ-123 [clean] - sample project
      └─ gws (branch: PROJ-123-backend)
    └─ PROJ-124 [dirty changes] - wip
      └─ gws (branch: PROJ-124-backend)
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

Manual (GitHub Releases):
- Download the archive for your OS/arch
- Extract and place `gws` on your PATH

For details and other options, see `docs/guides/INSTALL.md`.

## Quickstart (5 minutes)

### 1) Initialize the root

```bash
gws init
```

This creates `GWS_ROOT` with the standard layout and a starter `templates.yaml`.

Root resolution order:
1) `--root <path>`
2) `GWS_ROOT` environment variable
3) `~/gws` (default)

Default layout example:

```
~/gws/
  bare/        # bare repo store (shared Git objects)
  workspaces/  # task worktrees (one folder per workspace id)
  templates.yaml
```

### 2) Fetch repos (bare store)

```bash
gws repo get git@github.com:org/backend.git
```

### 3) Create a workspace

```bash
gws create --repo git@github.com:org/backend.git
```

You'll be prompted for a workspace id (e.g. `PROJ-123`).

Or run `gws create` with no args to pick a mode and fill inputs interactively.

### 4) Work and clean up

List workspaces:

```bash
gws ls
```

Open a workspace (prompts if omitted):

```bash
gws open PROJ-123
```

This launches an interactive subshell at the workspace root (parent cwd unchanged) and
prefixes the prompt with `[gws:<WORKSPACE_ID>]`.

Remove a workspace with guardrails (prompts if omitted):

```bash
gws rm PROJ-123
```

## Help and docs

- `docs/README.md` for documentation index
- `docs/spec/README.md` for specs index and status
- `docs/spec/commands/` for per-command specs (create/add/rm/etc.)
- `docs/spec/core/TEMPLATES.md` for template format
- `docs/spec/core/DIRECTORY_LAYOUT.md` for the file layout
- `docs/spec/ui/UI.md` for output conventions
- `docs/concepts/CONCEPT.md` for the background and motivation

## Maintainer

- @tasuku43
