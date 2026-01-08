---
title: "gws review"
status: implemented
---

## Synopsis
`gws review <PR URL>`

## Intent
Create a review-focused workspace for a GitHub pull request, checking out the PR head branch in a dedicated worktree.

## Behavior
- Accepts GitHub PR URLs only (e.g., `https://github.com/owner/repo/pull/123`); rejects other hosts or malformed paths.
- Uses `gh api` to fetch PR metadata (requires authenticated GitHub CLI): PR number, head ref, and repositories.
- Rejects forked PRs (head repo must match base repo).
- Selects the repo URL based on `defaultRepoProtocol` (SSH preferred, HTTPS fallback).
- Workspace ID is `REVIEW-PR-<number>`; errors if it already exists.
- Ensures the repo store exists, prompting to run `gws repo get` if missing (unless `--no-prompt`, which fails instead).
- Fetches the PR head ref into the bare store: `git fetch origin <head_ref>`.
- Adds a worktree under `<root>/workspaces/REVIEW-PR-<number>/<alias>` where:
  - Branch is `<head_ref>`.
  - Base ref is `refs/remotes/origin/<head_ref>`.
- Shows a summary of the workspace and worktree after creation.

## Success Criteria
- New workspace `REVIEW-PR-<number>` exists with a worktree checked out to the PR head branch.

## Failure Modes
- Invalid or unsupported PR URL; non-GitHub host.
- Fork PR detected.
- Missing or unauthenticated `gh` CLI.
- Repo store missing and user declines/forbids `repo get`.
- Git errors fetching the PR head or creating the worktree.
