# Manage Git worktrees declaratively with YAML (plan/apply + guardrails) — gion

**TL;DR:** I built **gion**, a CLI to manage **Git worktrees** declaratively.  
Write your desired state in `gion.yaml`, then **plan/apply** to create/update/cleanup in bulk **with guardrails**.  
It also supports **task-scoped workspaces** (grouping multiple repos for a **monorepo-like workflow**) and fast navigation via **giongo**.

Repo / docs:
- https://github.com/tasuku43/gion

---

## Why I built this

After I started using AI coding agents more often, I ended up doing more parallel work. That naturally led me to `git worktree`.

But the more parallel tasks I had, the more I kept asking myself:

- Where should I create the next worktree?
- Navigating across many worktrees is annoying.
- Can I delete this safely? (Or will I accidentally lose local changes?)

There are already tools around worktrees, but I wanted two things in particular:

1) **Guardrails during cleanup**  
   I wanted to prevent “I didn’t mean to delete it, but it was deletable anyway.”

2) **Task-scoped workspaces**  
   I wanted a “box per task” that can contain **one or more worktrees**, possibly across **multiple repositories**, so I can run my coding agent from the workspace root — and then delete the whole box when done.

That’s what **gion** is.

![Demo: gion reconciles workspaces via plan/apply (create/update/delete)](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/gdqcs7okhzr8ui9o7ln8.gif)

---

## Overview

The core workflow is **Create / Move / Cleanup**:

- **Create**: `gion manifest add` → `gion apply` (review the plan, then apply)
- **Move**: use `giongo` to search and jump
- **Cleanup**: `gion manifest gc` / `gion manifest rm`


Tip: `gion manifest` can be shortened to `gion m` or `gion man`.

---

## How it works: `gion.yaml` and the `manifest` subcommands

The center of gion is **`gion.yaml`**.

- `gion.yaml` is the **source of truth** (“desired state” / inventory).
- `gion manifest ...` is an **entry point** to update that YAML (you can also edit it directly).
- After any update (via command or direct editing), you run **`gion apply`** to reconcile the filesystem with the desired state.
- `gion apply` computes a **plan**, shows the diff, and then applies it when you confirm.

### Terminology: worktree vs workspace

- A **Git worktree** is a working directory that checks out a branch (or a specific commit).
- A **workspace** (in gion’s terms) is a **task-scoped directory** that can contain **one or more worktrees** — potentially from multiple repos.

A typical layout looks like this:

```text
GION_ROOT/
├─ gion.yaml           # desired state (inventory)
├─ bare/               # shared bare repo store
└─ workspaces/         # task-scoped workspaces
   ├─ PROJ-123/        # workspace_id (task)
   │  ├─ backend/      # worktree (repo: backend)
   │  ├─ frontend/
   │  └─ docs/
   └─ PROJ-456/
      └─ backend/
```

---

## Create: review a plan, then create in bulk

There are two ways to declare “create this workspace”:

- Use `gion manifest add`, or
- Edit `gion.yaml` directly

In both cases, you’re declaring a desired state first, and then reconciling with **`gion apply`**:
1) compute the plan
2) show the diff
3) apply when you confirm

### Four creation modes

`gion manifest add` supports four entry paths:

- `repo`
- `issue`
- `review`
- `preset`

The entry point differs, but the destination is the same: everything ends up in `gion.yaml` as desired state.

![Screenshot: gion manifest add modes (repo/issue/review/preset)](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/qeedcq6ufp7q11726t1c.png)

#### `issue` / `review`: queue up many, then apply once

If you use `issue` / `review`, you’ll need the `gh` CLI (GitHub-only).

The intended flow is:
- select multiple Issues / PRs
- add them to `gion.yaml`
- run `gion apply` once
- review the plan, then apply

This is great when you want to spin up many review/issue worktrees quickly.

![Demo: select multiple Issues/PRs and create worktrees in one apply](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/tfe4u3a1iz9uppv1kw4u.gif)

#### `repo`: create a single workspace quickly

If you just want the fastest “create one workspace”, `repo` is the simplest:
- choose a repo
- set a workspace id
- confirm the plan
- apply

![Screenshot: repo mode (create one workspace quickly)](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/mokwx7dosfgl3rtj7im7.png)

#### `preset`: group multiple repos under one task workspace (monorepo-like workflow)

A workspace is a “task box”, so it’s common to want something like:

- backend + frontend + docs
- backend + infra
- or any other multi-repo set

With **presets**, you define the set once, and then reuse it to create workspaces repeatedly.

![Screenshot: create a preset (define a reusable repo set)](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/57acy19xsjfgt4do8179.png)

![Screenshot: create a workspace from a preset (multi-repo)](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/7aqcmz5e0etprx10d46m.png)

#### Direct YAML editing

Direct editing is useful when you want to:
- adjust branch names in bulk
- create/remove multiple workspaces at once
- refactor/reorganize existing definitions

After editing, run:

```bash
gion apply
```

You’ll get a plan showing **creates / deletes / updates** together, so you can calmly review what will happen before applying.

![Screenshot: plan shows create/delete/update changes together](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/psperyjrfddicopxrbba.png)

---

## Move: search and jump to a workspace/worktree

Once worktrees start to accumulate, you waste time thinking:

> “Where was I working on that?”

For navigation, use **`giongo`** (installed alongside `gion` via Homebrew/mise).

- `giongo` does **not** change any state.
- It lists both **workspaces** and **worktrees**, and lets you filter and select.
- Select a **workspace** → jump to the workspace directory
- Select a **worktree** → jump to that repo’s working directory

![Demo: giongo lists and jumps to workspaces/worktrees](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/fi7w0edfr2ktvxpkvono.gif)

### Shell integration: make `giongo` change directories

If you want selection → `cd`, you need a small shell integration. The README includes options, but the simplest is:

```bash
eval "$(giongo init)"
```

(That prints a function you can paste into `~/.zshrc` or `~/.bashrc` for permanent setup.)

---

## Cleanup: conservative `gc`, and guarded `rm`

As worktrees pile up, you eventually pause and ask:

> “Is it safe to delete this?”

gion splits cleanup into two commands:

- `gion manifest gc`: **automatic, conservative cleanup**
- `gion manifest rm`: **manual selection with guardrails**

### `gion manifest gc`: conservative auto-cleanup

`gion manifest gc` proposes candidates that are **highly likely safe** to delete.

For example:
- workspaces whose checked-out branches are already merged into the default branch can be reclaimed
- anything ambiguous (uncommitted / unpushed / unreadable state, etc.) is excluded by default
- even “created but no commits” workspaces are excluded so you don’t delete something by accident

In short: `gc` aims to keep false-positives very low.

![Screenshot: gc classifies safe-to-delete candidates conservatively](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/90lzkd6aubtu9p87o4z5.png)

### `gion manifest rm`: manual removal, with guardrails

`gion manifest rm` is for cases where **a human decides** what to delete.

It supports interactive selection and then a final confirmation. During selection, each workspace gets lightweight tags like:

- `[dirty]`
- `[unpushed]`
- `[diverged]`
- `[unknown]`

What those mean:

- **dirty**: working tree has uncommitted changes (including untracked files or conflicts)
- **unpushed**: your local branch is ahead of upstream (has commits not pushed)
- **diverged**: local and upstream have both advanced (ahead and behind)
- **unknown**: cannot be determined (no upstream, detached HEAD, git error, etc.)

When you delete something risky (e.g., `dirty`), the **plan** will show a risk summary and (for dirty cases) the changed files, so you can sanity-check quickly before you confirm.

![Demo: rm shows risk and changed files in the plan before deletion](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/2jdlxkz8lggmv0l5kk8n.gif)

### Cleanup via direct YAML editing

You can also remove workspaces by editing `gion.yaml` directly.

Even then, `gion apply` will:
- clearly show what will be removed in the plan
- ask for a confirmation before destructive changes

So if you get nervous, you can just answer `n` and stop.

---

## Closing

Installation (Homebrew / mise) and full usage examples are in the GitHub README. If you’ve ever felt the pain of worktree sprawl — especially in multi-repo tasks — I’d love for you to try **gion** and share feedback.

- https://github.com/tasuku43/gion
- https://tasuku43.github.io/gion/
