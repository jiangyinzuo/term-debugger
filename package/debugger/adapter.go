package debugger

import (
	"io"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
)

type ProcessChildOutputFunc func(string, texteditor.TextEditor)

type TermDebugger interface {
	Name() string
	EndPromt() string
	ProcessChildStdout(string, texteditor.TextEditor)
	ProcessChildStderr(string, texteditor.TextEditor)
	ProcessUserInput(io.WriteCloser, string)
}
