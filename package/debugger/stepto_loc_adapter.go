package debugger

import (
	"regexp"
	"strconv"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
	log "github.com/sirupsen/logrus"
)

type steptoLocAdapter struct {
	baseDebugger

	steptoLocRegex *regexp.Regexp
}

func newSteptoAdapter(
	editorState *texteditor.EditorState,
	name, endPromt string,
	breakRegex, steptoLocRegex *regexp.Regexp) *steptoLocAdapter {
	return &steptoLocAdapter{
		baseDebugger: baseDebugger{
			editorState: editorState,
			name:        name,
			endPromt:    endPromt,
			breakRegex:  breakRegex,
		},
		steptoLocRegex: steptoLocRegex,
	}
}

func (s *steptoLocAdapter) ProcessChildOutput(output string) {
	switch output {
	case s.EndPromt():
		s.processPromptEnd()
	default:
		if s.maybeStepToLoc(output) || s.maybeAddBreakpoint(output) {
		}
	}
	s.baseDebugger.printChildOutput(output)

}

func (s *steptoLocAdapter) maybeStepToLoc(output string) bool {
	if matches := s.steptoLocRegex.FindStringSubmatch(output); matches != nil {
		filename := matches[1]
		line, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		log.Debugf("matches: %s:%d\n", filename, line)
		s.editorState.SetCurrentFileName(filename)
		s.editorState.SetCurrentLine(line)
		return true
	}
	return false
}
