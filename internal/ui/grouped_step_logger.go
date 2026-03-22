package ui

import "github.com/tasuku43/gion/internal/infra/output"

type groupedStepLogger struct {
	renderer     *Renderer
	current      *groupedStep
	currentGroup *groupedStepNode
}

type groupedStep struct {
	title string
	nodes []*groupedStepNode
}

type groupedStepNode struct {
	title string
}

func NewGroupedStepLogger(renderer *Renderer) *groupedStepLogger {
	return &groupedStepLogger{renderer: renderer}
}

func (l *groupedStepLogger) Step(text string) {
	l.Flush()
	l.current = &groupedStep{title: text}
}

func (l *groupedStepLogger) BeginGroup(text string) {
	if l.renderer == nil {
		return
	}
	if l.current == nil {
		l.renderer.StepLog(text)
		return
	}
	node := &groupedStepNode{title: text}
	l.current.nodes = append(l.current.nodes, node)
	l.currentGroup = node
}

func (l *groupedStepLogger) EndGroup() {
	l.currentGroup = nil
}

func (l *groupedStepLogger) Log(text string) {
	if l.renderer == nil {
		return
	}
	if l.current == nil {
		l.renderer.StepLog(text)
		return
	}
	if l.currentGroup == nil {
		l.current.nodes = append(l.current.nodes, &groupedStepNode{title: text})
	}
}

func (l *groupedStepLogger) LogOutput(text string) {
	return
}

func (l *groupedStepLogger) Flush() {
	if l.renderer == nil || l.current == nil {
		return
	}

	l.renderer.Step(l.current.title)
	for i, node := range l.current.nodes {
		prefix := output.TreeBranchMid
		if i == len(l.current.nodes)-1 {
			prefix = output.TreeBranchLast
		}
		l.renderer.TreeLine(
			l.renderer.MutedText(output.Indent+prefix),
			l.renderer.MutedText(node.title),
		)
	}

	l.current = nil
	l.currentGroup = nil
}
