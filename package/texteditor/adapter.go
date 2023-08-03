package texteditor

type TextEditor interface {
	StepTo(filename string, line int)
	SetBreakpoint(brkType BrkType, filename string, line int)
}
