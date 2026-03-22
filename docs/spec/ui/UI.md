---
title: "UI Specification"
status: implemented
---

# UI Specification (Bubble Tea + Bubbles + Lip Gloss)

## Goals
- Provide a quiet, consistent terminal UI
- Keep interactive and non-interactive output aligned under one information architecture
- Make prompt flows feel related to `kra` without forcing identical behavior

## Scope
- Interactive (TTY): Bubble Tea + Bubbles + Lip Gloss
- Non-interactive (non-TTY): plain text with the same layout rules

## Layout (common)

`gion` uses two related section families.

- Prompt flows use `Context` and `Step`
- Reconcile flows use `Info`, `Plan`, `Step` (optional confirm), `Apply`, and `Result`

Prompt-oriented layout:

```
Context
  <confirmed value>
  <pending muted value>

Info
  <warnings / blocked / auxiliary meta>

Step
  <active selector or question>
  <filter/input/footer>

Result
  <result line>

Suggestion
  <next command>
```

Reconcile-oriented layout:

```
Context
  <optional prelude from manifest mutation commands>

Info
  <manifest updated / warning lines>

Plan
  <plan line>
  <plan line>

Step
  <optional confirmation prompt>

Apply
  <execution summary>
  <workspace -> repo tree>

Result
  <result line>
```

Rules:
- Indent: 2 spaces
- 1 blank line between sections
- Section order is fixed within the active flow:
  - Prompt flow: `Context -> Info -> Step -> Result -> Suggestion`
  - Reconcile flow: `Context -> Info -> Plan -> Step -> Apply -> Result`
- Not every screen renders every section
- No success banner; success is implied in `Result`
- Long lines should wrap to terminal width; continuation lines keep the same text indent
  - Exception: interactive selector rows should truncate with `...` when stable cursor-to-row mapping matters
  - When truncation would hide the user-entered value at the end of a long label, prefer a 2-line layout: label on one line, then an indented detail line such as `branch: <value>`

Notes:
- `gion plan` renders `Info` (optional) then `Plan`
- `gion apply` renders `Info` (optional) -> `Plan` -> `Step` (interactive confirm only) -> `Apply` -> `Result`
- In normal output, `Apply` should summarize work at the workspace / repo level rather than dump raw git commands
- After confirmation completes, the `Step` prompt should not remain in the final visible transcript

## Prefix & Indentation
- Default prefix token: `•`
- Context / Info / Step / Result lines use `2 spaces + prefix + space` unless they are already formatted raw tree lines
- Tree/list indentation must use shared tokens from `internal/infra/output`

Prefix coloring:
- Info / Result / tree connectors: muted gray
- Prompts and active labels: accent
- Warnings: yellow
- Errors: red

Example:

```
Apply
  • create workspace PROJ-123
    └─ backend (branch: PROJ-123)
```

## Command execution logs
- Raw command logs should not appear in normal `Apply` output
- If commands are shown at all, render them in muted color
- Debug logging belongs in files when `--debug` is provided; do not add a visible `Debug` section

## Colors
- success: green
- warn: yellow
- error: red
- muted / log / connectors: low-contrast gray
- accent / active labels: cyan

## Components
- Text input: Bubbles `textinput`
- Selector / confirm: Bubble Tea models backed by shared rendering helpers
- Help line: muted, minimal
- Long selection lists should scroll so `Context` and `Step` remain visible

Selector-specific interaction and visual rules are documented in
[UI-SELECTOR.md](/Users/tasuku43/gionroot/workspaces/gion/gion/docs/spec/ui/UI-SELECTOR.md).

## Workspace picker line format

When rendering workspace candidates (for example in `gion manifest rm`), use:

`<WORKSPACE_ID>[<status>] - <description>`

Rules:
- `<description>` is optional
- Repo details should not be shown in the picker; deep review belongs in `Plan`
- `<status>` is omitted for clean workspaces
- Priority when multiple conditions apply:
  - `unknown` > `dirty` > `diverged` > `unpushed`

Color guidance:
- `<WORKSPACE_ID>`: default text color
- `[unpushed]` / `[diverged]`: warn
- `[dirty]` / `[unknown]`: error
- description: muted

## Examples

### `gion manifest add` (mode picker)

```
Context
  • mode:

Step
  • mode:
  > ○ repo - 1 repo only
    ○ issue - From an issue (multi-select, GitHub only)
    ○ review - From a review request (multi-select, GitHub only)
    ○ preset - From preset

  filter: s
  ↑↓ move  space/enter select  type filter  esc cancel
```

### `gion manifest add --preset` (interactive)

```
Context
  • mode: preset
  • preset: helmfiles
  • workspace id: PROJ-123

Info
  • manifest: updated gion.yaml

Plan
  • + add workspace PROJ-123
    └─ helmfiles (branch: PROJ-123)
       repo: github.com/org/helmfiles

Step
  • Apply changes? (default: No) (y/n)

  input: y
  enter confirm  esc cancel

Apply
  • create workspace PROJ-123
    └─ helmfiles (branch: PROJ-123)

Result
  • applied: add=1 update=0 remove=0
  • gion.yaml rewritten
```

### `gion manifest add --review` (interactive multi-select)

```
Context
  • mode: review
  • repo: git@github.com:org/repo.git
  • pull request:

Step
  • pull request:
    ○ #101 fix example
  > ● #102 align prompt UI

  filter:
  selected: 1/2  ↑↓ move  space toggle  enter apply  type filter  esc cancel
```

### `gion manifest add --preset` (`--no-apply`)

```
Context
  • mode: preset
  • preset: app
  • workspace id: PROJ-123
  • repo #1 (git@github.com:org/repo.git)
    └─ branch: PROJ-123

Result
  • updated gion.yaml

Suggestion
  gion apply
```

### `gion apply`

```
Plan
  • - remove workspace PROJ-099
    └─ backend (branch: PROJ-099)
       sync: upstream=origin/main ahead=1 behind=0

Step
  • Apply destructive changes? (default: No) (y/n)

  input: y
  enter confirm  esc cancel

Apply
  • remove workspace PROJ-099
    └─ backend (branch: PROJ-099)
```

## Notes
- All prompts and labels are English
- `Info` is optional and may include warnings, blocked items, or derived metadata
- `Suggestion` is optional and shown only on TTY with colors enabled

## Implementation contract
- CLI output must use `ui.Renderer` or `internal/infra/output` helpers
- Prompt models must use `Frame` rather than writing ad-hoc section headers
- Result lines must be rendered via `Bullet()`-style helpers to preserve prefix consistency
