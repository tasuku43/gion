---
title: "UI Architecture"
status: implemented
---

# UI Architecture

This document describes the implementation structure that keeps the terminal UI aligned with [UI.md](/Users/tasuku43/gionroot/workspaces/gion/gion/docs/spec/ui/UI.md).

## Goals
- Guarantee the fixed section order for the active flow
- Avoid per-command bespoke rendering
- Keep prompt state and reconcile summaries aligned with one shared contract
- Keep selector implementations aligned with [UI-SELECTOR.md](/Users/tasuku43/gionroot/workspaces/gion/gion/docs/spec/ui/UI-SELECTOR.md)

## Components

### Frame (`internal/ui/frame.go`)

Responsibilities:
- Centralize `Context / Info / Step / Result / Suggestion` ordering
- Hold already-classified lines and render them with consistent spacing

Usage:
- `SetContextPrompt(...)` / `AppendContextRaw(...)` for confirmed and pending context
- `SetStepPrompt(...)` / `AppendStepRaw(...)` for the active selector, text input, or confirmation
- `SetInfo(...)` / `AppendInfoRaw(...)` for warnings, blocked items, and auxiliary metadata

Key points:
- `Frame` owns prompt-screen structure
- Prompt models should update content only, not section ordering
- `Plan` / `Apply` remain CLI-rendered sections outside `Frame`

### Renderer (`internal/ui/renderer.go`)

Responsibilities:
- Low-level rendering of headers, bullets, raw lines, and tree indentation
- Shared by both `Frame` and CLI output paths

### Prompt Models (`internal/ui/prompt.go`)

Responsibilities:
- Own input / selection state transitions and validation
- Render prompt flows through `Frame`
- Centralize selector behavior instead of duplicating it in command handlers

## Implementation Rules
- Do not print UI output directly with `fmt.Fprintf/Printf/Println` in UI paths
- Consolidate confirmed state into `Context`, active interaction into `Step`, and auxiliary state into `Info`
- Do not invent ad-hoc sections for selected items or picker-specific summaries
- Do not use AltScreen

## Applying to Existing Flows
- Prefer one `Frame` that updates across multiple prompt steps
- Use `AppendContextRaw(...)` for pending context lines
- Use `AppendStepRaw(...)` for selector rows, `filter:` / `input:` lines, and footers
- Treat selector copy, assist-line wording, and confirm transitions as shared UI behavior
