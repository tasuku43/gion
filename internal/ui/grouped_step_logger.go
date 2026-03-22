package ui

import "github.com/tasuku43/gion/internal/infra/output"

type groupedStepLogger struct {
	renderer   *Renderer
	current    *groupedStep
	currentLog *groupedStepLog
}

type groupedStep struct {
	title string
	logs  []groupedStepLog
}

type groupedStepLog struct {
	text    string
	outputs []string
}

func NewGroupedStepLogger(renderer *Renderer) *groupedStepLogger {
	return &groupedStepLogger{renderer: renderer}
}

func (l *groupedStepLogger) Step(text string) {
	l.Flush()
	l.current = &groupedStep{title: text}
}

func (l *groupedStepLogger) Log(text string) {
	if l.renderer == nil {
		return
	}
	if l.current == nil {
		l.renderer.StepLog(text)
		return
	}
	l.current.logs = append(l.current.logs, groupedStepLog{text: text})
	l.currentLog = &l.current.logs[len(l.current.logs)-1]
}

func (l *groupedStepLogger) LogOutput(text string) {
	if l.renderer == nil {
		return
	}
	if l.current == nil || l.currentLog == nil {
		l.renderer.StepLogOutput(text)
		return
	}
	l.currentLog.outputs = append(l.currentLog.outputs, text)
}

func (l *groupedStepLogger) Flush() {
	if l.renderer == nil || l.current == nil {
		return
	}

	l.renderer.Step(l.current.title)
	for i, entry := range l.current.logs {
		prefix := output.TreeBranchMid
		if i == len(l.current.logs)-1 {
			prefix = output.TreeBranchLast
		}
		l.renderer.TreeLine(
			l.renderer.MutedText(output.Indent+prefix),
			l.renderer.MutedText(entry.text),
		)
		detailPrefix := l.renderer.MutedText(output.Indent + output.DetailTreePrefix(i == len(l.current.logs)-1))
		for _, line := range entry.outputs {
			l.renderer.TreeLine(detailPrefix, l.renderer.MutedText(line))
		}
	}

	l.current = nil
	l.currentLog = nil
}
