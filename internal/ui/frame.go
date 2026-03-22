package ui

import (
	"io"
	"strings"
)

// Frame renders a single-screen layout with fixed section order.
type Frame struct {
	Context    []frameLine
	Info       []frameLine
	Step       []frameLine
	Result     []frameLine
	Suggestion []frameLine

	theme    Theme
	useColor bool

	NoBlankAfterInfo bool
}

func NewFrame(theme Theme, useColor bool) *Frame {
	return &Frame{theme: theme, useColor: useColor}
}

func (f *Frame) SetContext(lines ...string) {
	f.Context = copyLines(lines, lineBullet)
}

func (f *Frame) SetContextPrompt(lines ...string) {
	f.Context = copyLines(lines, linePrompt)
}

func (f *Frame) AppendContextPrompt(lines ...string) {
	f.Context = append(f.Context, copyLines(lines, linePrompt)...)
}

func (f *Frame) SetContextRaw(lines ...string) {
	f.Context = copyRawLines(lines)
}

func (f *Frame) AppendContextRaw(lines ...string) {
	f.Context = append(f.Context, copyRawLines(lines)...)
}

func (f *Frame) SetStep(lines ...string) {
	f.Step = copyLines(lines, lineBullet)
}

func (f *Frame) SetStepPrompt(lines ...string) {
	f.Step = copyLines(lines, linePrompt)
}

func (f *Frame) AppendStep(lines ...string) {
	f.Step = append(f.Step, copyLines(lines, lineBullet)...)
}

func (f *Frame) AppendStepPrompt(lines ...string) {
	f.Step = append(f.Step, copyLines(lines, linePrompt)...)
}

func (f *Frame) SetStepRaw(lines ...string) {
	f.Step = copyRawLines(lines)
}

func (f *Frame) AppendStepRaw(lines ...string) {
	f.Step = append(f.Step, copyRawLines(lines)...)
}

func (f *Frame) SetInfo(lines ...string) {
	f.Info = copyLines(lines, lineBullet)
}

func (f *Frame) SetInfoPrompt(lines ...string) {
	f.Info = copyLines(lines, linePrompt)
}

func (f *Frame) AppendInfo(lines ...string) {
	f.Info = append(f.Info, copyLines(lines, lineBullet)...)
}

func (f *Frame) AppendInfoPrompt(lines ...string) {
	f.Info = append(f.Info, copyLines(lines, linePrompt)...)
}

func (f *Frame) SetInfoRaw(lines ...string) {
	f.Info = copyRawLines(lines)
}

func (f *Frame) AppendInfoRaw(lines ...string) {
	f.Info = append(f.Info, copyRawLines(lines)...)
}

func (f *Frame) SetResult(lines ...string) {
	f.Result = copyLines(lines, lineBullet)
}

func (f *Frame) AppendResult(lines ...string) {
	f.Result = append(f.Result, copyLines(lines, lineBullet)...)
}

func (f *Frame) SetSuggestion(lines ...string) {
	f.Suggestion = copyLines(lines, lineBullet)
}

func (f *Frame) Render() string {
	var b strings.Builder
	_, _ = f.WriteTo(&b)
	return b.String()
}

func (f *Frame) WriteTo(w io.Writer) (int64, error) {
	cw := &countingWriter{w: w}
	r := NewRenderer(cw, f.theme, f.useColor)

	if len(f.Context) > 0 {
		r.Section("Context")
		for _, line := range f.Context {
			renderLine(r, line)
		}
		if len(f.Info) > 0 || len(f.Step) > 0 || len(f.Result) > 0 || len(f.Suggestion) > 0 {
			r.Blank()
		}
	}

	if len(f.Info) > 0 {
		r.Section("Info")
		for _, line := range f.Info {
			renderLine(r, line)
		}
		if !f.NoBlankAfterInfo && (len(f.Step) > 0 || len(f.Result) > 0 || len(f.Suggestion) > 0) {
			r.Blank()
		}
	}

	if len(f.Step) > 0 {
		r.Section("Step")
		for _, line := range f.Step {
			renderLine(r, line)
		}
		if len(f.Result) > 0 || len(f.Suggestion) > 0 {
			r.Blank()
		}
	}

	if len(f.Result) > 0 {
		r.Section("Result")
		for _, line := range f.Result {
			renderLine(r, line)
		}
		if len(f.Suggestion) > 0 && f.useColor {
			r.Blank()
		}
	}

	if len(f.Suggestion) > 0 && f.useColor {
		r.Section("Suggestion")
		for _, line := range f.Suggestion {
			renderLine(r, line)
		}
	}

	return cw.n, cw.err
}

type frameLine struct {
	text string
	kind frameLineKind
}

type frameLineKind int

const (
	lineBullet frameLineKind = iota
	linePrompt
	lineStep
	lineRaw
)

type countingWriter struct {
	w   io.Writer
	n   int64
	err error
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	if cw.err != nil {
		return 0, cw.err
	}
	n, err := cw.w.Write(p)
	cw.n += int64(n)
	if err != nil {
		cw.err = err
	}
	return n, err
}

func renderLine(r *Renderer, line frameLine) {
	if line.kind == lineRaw && line.text == "" {
		r.LineRaw("")
		return
	}
	if strings.TrimSpace(line.text) == "" {
		return
	}
	switch line.kind {
	case lineRaw:
		r.LineRaw(line.text)
	case linePrompt:
		r.Prompt(line.text)
	case lineStep:
		r.Step(line.text)
	default:
		r.Bullet(line.text)
	}
}

func (f *Frame) SetInputs(lines ...string) {
	f.SetContext(lines...)
}

func (f *Frame) SetInputsPrompt(lines ...string) {
	f.SetContextPrompt(lines...)
}

func (f *Frame) AppendInputsPrompt(lines ...string) {
	f.AppendContextPrompt(lines...)
}

func (f *Frame) SetInputsRaw(lines ...string) {
	f.SetContextRaw(lines...)
}

func (f *Frame) AppendInputsRaw(lines ...string) {
	f.AppendContextRaw(lines...)
}

func copyLines(lines []string, kind frameLineKind) []frameLine {
	var out []frameLine
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		out = append(out, frameLine{text: trimmed, kind: kind})
	}
	return out
}

func copyRawLines(lines []string) []frameLine {
	var out []frameLine
	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\n")
		if trimmed == "" {
			out = append(out, frameLine{text: "", kind: lineRaw})
			continue
		}
		if strings.TrimSpace(trimmed) == "" {
			continue
		}
		out = append(out, frameLine{text: trimmed, kind: lineRaw})
	}
	return out
}
