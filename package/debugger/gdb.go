package debugger

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
)

const (
	StateInitial = iota
	StateIdle
	StateBreak
	StateDeleteAll
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

type GDBAdapter struct {
	editor texteditor.TextEditor
	state  int
}

func NewGDBAdapter(editor texteditor.TextEditor) *GDBAdapter {
	return &GDBAdapter{
		editor: editor,
		state:  StateInitial,
	}
}

func (g *GDBAdapter) Name() string {
	return "gdb"
}

func (g *GDBAdapter) EndPromt() string {
	return "(gdb) "
}

func (g *GDBAdapter) ProcessChildOutput(output string) {
	switch g.state {
	case StateInitial:
		g.maybeAddBreakpoint(output)
	case StateBreak:
		g.maybeAddBreakpoint(output)
		g.state = StateIdle
	}

	if output == " Delete all breakpoints? (y or n) " {
		g.state = StateDeleteAll
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
	case StateDeleteAll:
		if input == "y" {
			g.editor.DeleteAllBreakPoints()
		} else {
			g.state = StateIdle
		}
	}
	stdin.Write([]byte(input + "\n"))
}

func (g *GDBAdapter) maybeAddBreakpoint(output string) {
	if matches := breakRegex.FindStringSubmatch(output); matches != nil {
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
		}
		line, err := strconv.Atoi(matches[3])
		if err != nil {
			panic(err)
		}
		g.editor.AddBreakPoint(id, matches[2], line)
	}
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

func (g *GDBAdapter) parseBreakPointByLoc(matches []string) {
	filename := matches[2]
	line, err := strconv.Atoi(matches[3])
	if err != nil {
		panic(err)
	}
	fmt.Println(filename, line)
}

func (g *GDBAdapter) maybeInputDeleteBreakpoint(input string) bool {
	fmt.Println("Maybe!")
	if matches := clearLocRegex.FindStringSubmatch(input); matches != nil {
		g.parseBreakPointByLoc(matches)
		return true
	}
	if matches := clearLineRegex.FindStringSubmatch(input); matches != nil {
		line, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		fmt.Println("TODO", line)
		return true
	}
	if matches := deleteLocRegex.FindStringSubmatch(input); matches != nil {
		g.parseBreakPointByLoc(matches)
		return true
	}
	if matches := deleteIdRegex.FindStringSubmatch(input); matches != nil {
		id, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		fmt.Println(id)
		return true
	}
	fmt.Println("Oops!")
	return false
}
