---
title: "gion manifest preset rm"
status: implemented
aliases:
  - "gion manifest pre rm"
  - "gion manifest p rm"
---

## Synopsis
`gion manifest preset rm [<name> ...] [--no-prompt]`

## Intent
Remove preset entries from `gion.yaml`.

## Notes
- This command is inventory-only and does not run `gion apply`.

## Behavior
- Accepts zero or more preset names. When multiple names are provided, duplicates are removed while preserving first-seen order.
- Requires `gion.yaml` to exist (`gion init` completed). Missing file => error.
- With names provided:
  - Errors if any requested name does not exist; no changes are written.
  - Otherwise removes the listed presets and writes the file back via atomic tmp+rename.
- With no names provided and prompts allowed:
  - Opens a filterable list of existing preset names (case-insensitive substring match).
  - Multi-select is supported.
  - Cancel/empty selection exits 0 with no changes.
  - With `--no-prompt`, error instead of opening the selector.
- After removal, other presets remain unchanged; repo stores are untouched.
- Output uses the common sectioned layout from `docs/spec/ui/UI.md`. No `Plan`/`Apply` sections are used.

## Interactive selection UX (no args)
- Candidate list is the preset names in `gion.yaml`.
- Prompt behavior mirrors existing gion selection UI:
  - Shows a filterable list. Typing narrows candidates by substring match (case-insensitive). Optionally a lightweight fuzzy match is acceptable.
  - The first visible item is highlighted with `>`.
  - `Space` toggles the highlighted preset with `○/●` markers.
  - `Enter` applies the selected set.
  - The active query is shown in the bottom `filter:` line.

## Output examples

### Output: interactive (no args)
```
Context
  • preset:

Step
  • preset:
    ○ helpdesk
  > ● helpers

  filter: he
  selected: 1/2  ↑↓ move  space toggle  enter apply  type filter  esc cancel

Result
  • updated gion.yaml (removed 2 presets)
```

### Output: non-interactive (args)
```
Context
  • preset: helpdesk

Result
  • updated gion.yaml (removed 1 preset)
```

## Success Criteria
- `gion.yaml` no longer contains the removed preset entries.

## Failure Modes
- `gion.yaml` missing or unreadable.
- Preset name not found (when explicitly provided).
- Write/rename failure when persisting the updated file.
