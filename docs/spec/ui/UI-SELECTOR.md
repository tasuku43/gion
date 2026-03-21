---
title: "Selector UI"
status: proposed
---

# Selector UI

This document defines the `gion` selector contract for interactive pickers.
It intentionally aligns with `kra` on visual language and terminal layout, while
keeping `gion`'s flow-oriented interaction model.

## Goals

- Make `gion` and `kra` feel like the same family in terminal UI quality
- Preserve `gion`'s input-flow orientation for manifest and workspace creation flows
- Reduce per-picker drift by documenting one selector contract for `gion`

## Design stance

`gion` selectors are not a direct copy of `kra` selectors.

- `kra` favors in-list toggle interaction for repeated operational tasks
- `gion` favors staged input flows that move the user toward one apply action

Because of that, `gion` keeps its own selection model, but should reuse the same
visual grammar where possible:

- stable section layout
- consistent filter input placement
- consistent muted/accent/error semantics
- consistent focus visibility
- consistent key-hint language

## Selector families

`gion` currently uses two selector families.

### Single-select

Used when the next step needs exactly one value.

Examples:
- mode selection in create flows
- repo selection in create flows
- workspace selection in simple pickers
- `giongo` workspace/repo selection

Interaction:
- `Up` / `Down`: move cursor
- `Enter`: confirm focused row
- `Esc` / `Ctrl+C`: cancel
- typing: update filter query
- `Backspace` / `Delete`: edit filter query

### Multi-select

Used when the user builds a set and then proceeds to the next step or apply.

Examples:
- review PR selection
- issue selection
- preset repo selection
- manifest workspace removal selection

Interaction:
- `Up` / `Down`: move cursor
- `Enter`: add the focused item into the selected set
- `Ctrl+D`: finish selection
- `done` + `Enter`: finish selection
- `Esc` / `Ctrl+C`: cancel
- typing: update filter query
- `Backspace` / `Delete`: edit filter query

Important distinction from `kra`:
- `gion` multi-select is append-and-progress, not toggle-in-place
- selected items may move out of the candidate list and into `Info`

This is intentional and matches `gion`'s staged workflow design.

## Shared layout contract

Selectors must follow [UI.md](/Users/tasuku43/gionroot/workspaces/gion/gion/docs/spec/ui/UI.md).

Rules:
- keep selector content inside the shared `Frame`
- selector input lives in `Inputs`
- selected items, blocked items, and validation messages live in `Info`
- do not introduce ad-hoc section names outside the shared section system
- do not use AltScreen

For selector rendering:
- the filter input line stays at the top of `Inputs`
- the candidate list is rendered immediately below the input line
- helper text such as finish hints is rendered after the candidate list
- selected items are rendered in `Info`
- blocked items, when present, are rendered after selected items in `Info`

## Visual language

To keep the family resemblance with `kra`, selectors should converge on the
same terminal language even when the interaction differs.

### Focus

- the focused row must be visually obvious
- in color terminals, focused text may use bold and/or a subtle emphasis style
- in no-color terminals, focus must remain understandable from row position alone

### Selection state

`gion` does not need to adopt `kra`'s `○/●` row markers everywhere.

Instead:
- candidate rows may remain compact and flow-oriented
- the selected set must be visible as a separate list or tree in `Info`
- the selected set should use accent for its heading and muted rendering for tree structure

### Text hierarchy

- primary information: selected item labels, focused row content
- muted information: descriptions, tree connectors, helper text
- error information: validation lines and blocking messages
- accent information: selector labels and selected-set headings

## Copy and hints

Selector copy should be short and explicit.

Recommended wording:
- `finish: Ctrl+D or type "done"` for append-style multi-select
- `select at least one <item>` for empty-confirm validation
- `no matches` when filter results are empty

`gion` may keep its own finish hint instead of adopting `kra`'s footer contract,
because the interaction model is different.

However, wording should still feel related:
- use imperative verbs
- keep hint phrases short
- prefer consistent English labels across commands

## Search behavior

Filtering should feel consistent across selectors.

Minimum contract:
- case-insensitive
- direct typing updates the filter immediately
- `Backspace` and `Delete` modify the filter without entering a separate mode
- when filter results shrink, cursor must stay within visible range

Current `gion` matching may remain simple substring-based where already implemented.
If `gion` later adopts the same fuzzy ordered-subsequence matching as `kra`, that
should be treated as an enhancement, not a compatibility requirement for this spec.

## Multi-select behavior

For append-style multi-select in `gion`:

- `Enter` adds the focused item to the selected set
- after adding, the item may be removed from remaining candidates
- the filter may be cleared after addition when that helps repeated entry
- finishing requires at least one selected item
- selected items must remain reviewable before the next stage

Recommended UX rules:
- keep the selected set visible without scrolling away from the input
- keep finish instructions visible near the candidate list
- do not require remembering hidden controls
- avoid mixing selected and blocked states into the same visual block

## Architecture guidance

`gion` should still avoid per-command bespoke selectors.

Preferred direction:
- keep selector state machines in `internal/ui/prompt.go`
- reuse shared rendering helpers for candidate rows and selected trees
- keep `Frame` responsible for section ordering
- keep command handlers responsible for supplying data, not formatting

This does not require extracting a shared library with `kra`.
That question should be revisited only after the interaction and visual contracts
have stabilized independently in both projects.
