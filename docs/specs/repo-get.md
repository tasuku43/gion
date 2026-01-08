---
title: "gws repo get"
status: implemented
---

## Synopsis
`gws repo get <repo>`

## Intent
Create or normalize a bare repo store for a remote Git repository and ensure a matching working copy exists under `src/`.

## Behavior
- Accepts SSH or HTTPS Git URLs (e.g., `git@github.com:owner/repo.git` or `https://github.com/owner/repo.git`).
- Normalizes the repo spec to derive a stable repo key and store path (`<root>/bare/<host>/<owner>/<repo>.git`).
- If the store is missing, clones it as `--bare`.
- Normalizes the store:
  - Sets `remote.origin.fetch` to `+refs/heads/*:refs/remotes/origin/*`.
  - Detects the default branch from the remote and updates `refs/remotes/origin/HEAD` accordingly.
  - Runs `git fetch --prune` when the local store is stale.
  - Prunes local head refs that no longer exist remotely.
- Ensures a non-bare working copy exists at `<root>/src/<host>/<owner>/<repo>`:
  - Clones from the bare store when missing.
  - Updates `origin` URL to the original remote.
  - Does not fetch the working copy when it already exists (MVP behavior).

## Success Criteria
- Bare store exists, normalized, and up to date with the remote default branch.
- Working copy exists (cloned from the bare store) with `origin` pointing to the remote URL.

## Failure Modes
- Missing repo argument or invalid repo spec.
- Network or git errors during clone/fetch.
- Filesystem errors creating store or src paths.
