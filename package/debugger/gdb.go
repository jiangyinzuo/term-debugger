package debugger

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
)

const (
	listFlag uint64 = 1 << iota
)

const gdbLocRegex = `([a-zA-Z0-9_\-\.\/]+):(\d+)`

var gdbStepToCursorLocRegex = regexp.MustCompile(fmt.Sprintf(`at %s\n$`, gdbLocRegex))
var gdbStepToLineRegex = regexp.MustCompile(`^(\d+)\s+`)

// Example:
// Breakpoint 1 at 0x11e7: file hello.cpp, line 18.
var gdbBreakRegex = regexp.MustCompile(`^Breakpoint (\d+) at .+: file (.+), line (\d+)`)
var gdbListRegex = regexp.MustCompile(`^list|lis|li$`)

type GDBAdapter struct {
	baseDebugger

	flags uint64
}

func NewGDBAdapter(editorState *texteditor.EditorState) *GDBAdapter {
	return &GDBAdapter{
		baseDebugger: baseDebugger{
			editorState: editorState,
			name:        "gdb",
			endPromt:    "(gdb) ",
			breakRegex:  gdbBreakRegex,
		},
	}
}

func (g *GDBAdapter) ProcessChildOutput(output string) {
	switch output {
	case g.EndPromt():
		g.processPromptEnd()
	default:
		if g.maybeStepToLoc(output) || g.maybeAddBreakpoint(output) {
		}
	}
	g.baseDebugger.printChildOutput(output)
}

func (g *GDBAdapter) ProcessUserInput(stdin io.WriteCloser, input string) {
	trimmed := strings.TrimSpace(input)
	if gdbListRegex.MatchString(trimmed) {
		g.flags |= listFlag
	}
	g.baseDebugger.ProcessUserInput(stdin, input)
}

func (g *GDBAdapter) processPromptEnd() {
	g.flags = 0
	g.baseDebugger.processPromptEnd()
}

func (g *GDBAdapter) maybeStepToLoc(output string) bool {
	if g.flags&listFlag != 0 {
		return false
	}
	if matches := gdbStepToCursorLocRegex.FindStringSubmatch(output); matches != nil {
		filename := matches[1]
		line, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		g.editorState.SetCurrentFileName(filename)
		g.editorState.SetCurrentLine(line)
		return true
	}
	if matches := gdbStepToLineRegex.FindStringSubmatch(output); matches != nil {
		line, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
		}
		g.editorState.SetCurrentLine(line)
		return true
	}
	return false
}
