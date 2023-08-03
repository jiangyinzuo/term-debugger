package debugger

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
)

const (
	StateInitial = iota
	StateIdle
	StateBreak
)

var breakCommands = []string{"b", "br", "bre", "brea", "break"}

// Example:
//
//	Breakpoint 1 at 0x11e7: file hello.cpp, line 18.
var breakRegex = regexp.MustCompile(`^Breakpoint (\d+) at .+: file (.+), line (\d+)`)

const locRegex = `([a-zA-Z0-9_\-\.\/]+):(\d+)`

const clearCmdRegex = `(clear|cl)\s+`

var clearLocRegex = regexp.MustCompile(fmt.Sprintf("%s%s", clearCmdRegex, locRegex))
var clearLineRegex = regexp.MustCompile(fmt.Sprintf("%s%s", clearCmdRegex, `(\d+)`))

const deleteCmdRegex = `(delete|del|d)\s+`

var deleteLocRegex = regexp.MustCompile(fmt.Sprintf("%s%s", deleteCmdRegex, locRegex))
var deleteIdRegex = regexp.MustCompile(fmt.Sprintf("%s%s", deleteCmdRegex, `(\d+)`))

var stepToCursorLocRegex = regexp.MustCompile(fmt.Sprintf(`at %s\n$`, locRegex))
var stepToLineRegex = regexp.MustCompile(`^(\d+)\s+`)

type GDBAdapter struct {
	editorState *texteditor.EditorState
	state       int
}

func NewGDBAdapter(editorState *texteditor.EditorState) *GDBAdapter {
	return &GDBAdapter{
		editorState: editorState,
		state:       StateInitial,
	}
}

func (g *GDBAdapter) Name() string {
	return "gdb"
}

func (g *GDBAdapter) EndPromt() string {
	return "(gdb) "
}

func (g *GDBAdapter) ProcessChildOutput(output string) {
	switch output {
	case g.EndPromt():
		g.processPromptEnd()
	case "Delete all breakpoints? (y or n) [answered Y; input not from terminal]\n":
		log.Debugln("delete all breakpoints")
		g.editorState.RemoveAllBreakPoints()
	default:
		if g.maybeAddBreakpoint(output) ||
			g.maybeStepToLoc(output) {
		}
	}

	fmt.Print(output)
	os.Stdout.Sync()
}

func (g *GDBAdapter) ProcessUserInput(stdin io.WriteCloser, input string) {
	switch g.state {
	case StateInitial:
		fallthrough
	case StateIdle:
		if g.maybeInputAddBreakPoint(input) {
			break
		}
		if g.maybeInputDeleteBreakpoint(input) {
			break
		}
	}
	stdin.Write([]byte(input + "\n"))
}

func (g *GDBAdapter) processPromptEnd() {
	g.editorState.StepToCursor()
}

func (g *GDBAdapter) maybeStepToLoc(output string) bool {
	if matches := stepToCursorLocRegex.FindStringSubmatch(output); matches != nil {
		filename := matches[1]
		line, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		g.editorState.SetCurrentFileName(filename)
		g.editorState.SetCurrentLine(line)
		return true
	}
	if matches := stepToLineRegex.FindStringSubmatch(output); matches != nil {
		line, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
		}
		g.editorState.SetCurrentLine(line)
		return true
	}
	return false
}

func (g *GDBAdapter) maybeAddBreakpoint(output string) bool {
	if matches := breakRegex.FindStringSubmatch(output); matches != nil {
		log.Debugf("matches: %v\n", matches)
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
		}
		line, err := strconv.Atoi(matches[3])
		if err != nil {
			panic(err)
		}
		g.editorState.AddBreakpoint(id, matches[2], line)
		return true
	}
	return false
}

func (g *GDBAdapter) maybeInputAddBreakPoint(input string) bool {
	for _, command := range breakCommands {
		if strings.HasPrefix(input, command) {
			g.state = StateBreak
			return true
		}
	}
	return false
}

func (g *GDBAdapter) deleteBreakPointByLoc(matches []string) {
	filename := matches[2]
	line, err := strconv.Atoi(matches[3])
	if err != nil {
		panic(err)
	}
	g.editorState.RemoveBreakpointByLoc(filename, line)
}

func (g *GDBAdapter) maybeInputDeleteBreakpoint(input string) bool {
	if matches := clearLocRegex.FindStringSubmatch(input); matches != nil {
		g.deleteBreakPointByLoc(matches)
		return true
	}
	if matches := clearLineRegex.FindStringSubmatch(input); matches != nil {
		line, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		g.editorState.RemoveBreakpointByLine(line)
		return true
	}
	if matches := deleteLocRegex.FindStringSubmatch(input); matches != nil {
		g.deleteBreakPointByLoc(matches)
		return true
	}
	if matches := deleteIdRegex.FindStringSubmatch(input); matches != nil {
		id, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		g.editorState.RemoveBreakpointById(id)
		return true
	}
	return false
}
