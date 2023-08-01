package debugger

import (
	"io"
)

type ProcessChildOutputFunc func(string)

type TermDebugger interface {
	Name() string
	EndPromt() string
	ProcessChildOutput(string)
	ProcessUserInput(io.WriteCloser, string)
}
