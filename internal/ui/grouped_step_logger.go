package ui

import "github.com/tasuku43/gion/internal/infra/output"

type groupedStepLogger struct {
	renderer     *Renderer
	current      *groupedStep
	currentGroup *groupedStepNode
	currentNode  *groupedStepNode
}

type groupedStep struct {
	title string
	nodes []*groupedStepNode
}

type groupedStepNode struct {
	title    string
	outputs  []string
	children []*groupedStepNode
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
	l.currentNode = nil
}

func (l *groupedStepLogger) EndGroup() {
	l.currentGroup = nil
	l.currentNode = nil
}

func (l *groupedStepLogger) Log(text string) {
	if l.renderer == nil {
		return
	}
	if l.current == nil {
		l.renderer.StepLog(text)
		return
	}
	node := &groupedStepNode{title: text}
	if l.currentGroup != nil {
		l.currentGroup.children = append(l.currentGroup.children, node)
	} else {
		l.current.nodes = append(l.current.nodes, node)
	}
	l.currentNode = node
}

func (l *groupedStepLogger) LogOutput(text string) {
	if l.renderer == nil {
		return
	}
	if l.current == nil || l.currentNode == nil {
		l.renderer.StepLogOutput(text)
		return
	}
	l.currentNode.outputs = append(l.currentNode.outputs, text)
}

func (l *groupedStepLogger) Flush() {
	if l.renderer == nil || l.current == nil {
		return
	}

	l.renderer.Step(l.current.title)
	l.renderNodes(output.Indent, l.current.nodes)

	l.current = nil
	l.currentGroup = nil
	l.currentNode = nil
}

func (l *groupedStepLogger) renderNodes(baseIndent string, nodes []*groupedStepNode) {
	for i, node := range nodes {
		prefix := output.TreeBranchMid
		if i == len(nodes)-1 {
			prefix = output.TreeBranchLast
		}
		l.renderer.TreeLine(
			l.renderer.MutedText(baseIndent+prefix),
			l.renderer.MutedText(node.title),
		)

		childIndent := baseIndent + output.DetailTreePrefix(i == len(nodes)-1)
		if len(node.children) > 0 {
			l.renderNodes(childIndent, node.children)
		}
		for _, line := range node.outputs {
			l.renderer.TreeLine(
				l.renderer.MutedText(childIndent),
				l.renderer.MutedText(line),
			)
		}
	}
}
