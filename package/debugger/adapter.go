package debugger

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
	log "github.com/sirupsen/logrus"
)

type TermDebugger interface {
	Name() string
	EndPromt() string
	ProcessChildOutput(string)
	ProcessUserInput(io.WriteCloser, string)
}

type baseDebugger struct {
	editorState *texteditor.EditorState
	name        string
	endPromt    string

	breakRegex *regexp.Regexp
}

func (d *baseDebugger) Name() string {
	return d.name
}

func (d *baseDebugger) EndPromt() string {
	return d.endPromt
}

func (d *baseDebugger) ProcessUserInput(stdin io.WriteCloser, input string) {
	stdin.Write([]byte(input + "\n"))
}

func (d *baseDebugger) printChildOutput(output string) {
	fmt.Print(output)
	os.Stdout.Sync()
}

func (d *baseDebugger) processPromptEnd() {
	d.editorState.SyncCursor()
}

func (d *baseDebugger) maybeAddBreakpoint(output string) bool {
	if matches := d.breakRegex.FindStringSubmatch(output); matches != nil {
		_, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
		}
		line, err := strconv.Atoi(matches[3])
		if err != nil {
			panic(err)
		}
		d.editorState.SignBreakpoint(texteditor.BreakpointNormal, matches[2], line)
		log.Debugf("breakpoint added: %s:%d\n", matches[2], line)
		return true
	}
	return false
}
