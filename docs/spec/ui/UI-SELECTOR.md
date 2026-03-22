---
title: "Selector UI"
status: proposed
---

# Selector UI

This document defines the `gion` selector contract for interactive pickers.
It intentionally aligns with `kra` on visual language and terminal layout,
while keeping `gion`'s staged flow and summary-oriented apply behavior.

## Goals
- Make `gion` and `kra` feel like the same family in terminal UI quality
- Preserve `gion`'s staged flow for manifest and workspace creation
- Reduce per-picker drift by documenting one selector contract

## Design stance

`gion` does not need to copy `kra` interaction-for-interaction.

What should match:
- stable section layout
- `>`, `○/●`, muted connectors, and assist-line language
- consistent `filter:` / `input:` placement
- consistent focus, confirm, and error handling

What may differ:
- whether a flow proceeds to another step or applies immediately
- whether confirmation is visible only transiently
- how command-specific context is accumulated above the selector

## Selector families

### Single-select

Used when the next stage needs exactly one value.

Examples:
- mode selection in create flows
- repo selection in create flows
- workspace selection in simple pickers
- `giongo` workspace / repo selection

Interaction:
- `Up` / `Down`: move cursor
- `Space` / `Enter`: confirm focused row
- `Esc` / `Ctrl+C`: cancel
- typing: update filter query
- `Backspace` / `Delete`: edit filter query

### Multi-select

Used when the user builds a set and then applies or advances once.

Examples:
- review PR selection
- issue selection
- preset repo selection
- manifest workspace removal selection

Interaction:
- `Up` / `Down`: move cursor
- `Space`: toggle the focused item
- `Enter`: apply the current selected set
- `Esc` / `Ctrl+C`: cancel
- typing: update filter query
- `Backspace` / `Delete`: edit filter query

## Shared layout contract

Selectors must follow [UI.md](/Users/tasuku43/gionroot/workspaces/gion/gion/docs/spec/ui/UI.md).

Rules:
- keep selector content inside the shared `Frame`
- confirmed state lives in `Context`
- the current unresolved field also appears in `Context` as a muted pending line
- the active selector prompt, candidate list, and assist lines live in `Step`
- blocked items and validation messages live in `Info`
- do not introduce ad-hoc section names
- do not use AltScreen

For selector rendering:
- the active prompt line stays at the top of `Step`
- the candidate list is rendered immediately below the prompt line
- `filter:` and the footer are rendered after the candidate list
- confirm state may temporarily hide assist lines before the next stage

## Visual language

### Focus
- the focused row must be visually obvious
- no-color terminals must still show focus through row markers and placement

### Selection state
- use `○/●` markers for selector rows
- single-select shows the confirmed or focused choice with `●`
- multi-select keeps the selected set visible inline in the list
- do not move selected items into a separate `Info` block in normal selector flows

### Text hierarchy
- primary information: selected item labels, focused row content
- muted information: descriptions, connectors, helper text, pending context
- error information: validation lines and blocking messages
- accent information: active labels and selected markers

## Copy and hints

Recommended wording:
- `filter:` for selector filtering
- `input:` for text entry
- `select at least one <item>` for empty-confirm validation
- `no matches` when filtering returns nothing
- `↑↓ move  space/enter select  type filter  esc cancel` for single-select
- `selected: n/m  ↑↓ move  space toggle  enter apply  type filter  esc cancel` for multi-select

## Search behavior
- case-insensitive
- direct typing updates immediately
- `Backspace` and `Delete` edit the query in place
- when the result set shrinks, cursor must stay within visible range

Current matching may remain substring-based or fuzzy where already implemented.

## Architecture guidance
- keep selector state machines in `internal/ui/prompt.go`
- reuse shared rendering helpers for candidate rows, assist lines, and pending context lines
- keep `Frame` responsible for prompt section ordering
- keep command handlers responsible for supplying data, not formatting

This does not require a shared extracted library with `kra` yet.
