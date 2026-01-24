# giongo implementation tasks

## Status
- overall: implemented

## Scope
- Add a separate `giongo` binary focused on interactive workspace/worktree navigation.
- Reuse the existing picker UI and scanning logic from gion where possible.
- Keep gion's IaC command surface unchanged.

## Tasks
- [x] Spec and docs
- [x] CLI entrypoint (`cmd/giongo`)
- [x] Picker model extension (workspace + repo selectable rows)
- [x] Filesystem scanning + metadata description support
- [x] Search filtering (workspace + repo + details)
- [x] Selection output (absolute path, `--print`)
- [x] Non-TTY error behavior
- [x] GoReleaser config (builds + archives for `giongo`)
- [x] Tests (filtering, parent visibility, path resolution, non-TTY)
- [x] README/INSTALL updates

## Notes
- Update task statuses as work progresses.
