package texteditor

const (
	BreakpointNormal byte = 1
)

type TextEditor interface {
	StepTo(filename string, line int)
	SignBreakpoint(brkType byte, filename string, line int)
}
