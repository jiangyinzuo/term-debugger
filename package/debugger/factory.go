package debugger

import (
	"log"
	"regexp"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
)

// Example:
// Breakpoint 1 at /root/term-debugger/example/hello.py:5
var pdbBreakRegex = regexp.MustCompile(`^Breakpoint (\d+) at (.+):(\d+)`)

// Example:
// > /root/term-debugger/example/myplus.py(2)plus()
var pdbSteptoLocRegex = regexp.MustCompile(`^> (.+)\((\d+)\)`)

// Example:
// Breakpoint 1: where = hello`main + 27 at hello.cpp:8:16, address = 0x000000000000129c
var lldbBreakRegex = regexp.MustCompile(`^Breakpoint (\d+): where = .+ at ([^,]+):(\d+):\d+, address = `)

// Example:
// frame #0: 0x00005555555552e7 hello`main at hello.cpp:12:10
var lldbSteptoLocRegex = regexp.MustCompile(`^\s*frame #\d+: .+ at ([^:]+):(\d+):\d+`)

// Example:
// Breakpoint 1 set at 0x49b7ea for main.main() ./m.go:7
var dlvBreakRegex = regexp.MustCompile(`^Breakpoint (\d+) set at 0x[0-9a-f]+ for .+ ([^:]+):(\d+)`)

// Example:
// > main.main() ./m.go:11 (PC: 0x49b82f)
// > main.main() ./m.go:7 (hits goroutine(1):1 total:1) (PC: 0x49b7ea)
var dlvSteptoLocRegex = regexp.MustCompile(`^> .+ ([^:]+):(\d+) \(`)

func NewDebugger(debugTool *string) TermDebugger {
	textEditor := &texteditor.VimTapi{}
	editorState := texteditor.NewEditorState(textEditor)
	switch *debugTool {
	case "gdb":
		return NewGDBAdapter(editorState)
	case "pdb":
		return newSteptoAdapter(editorState, "pdb", "(Pdb) ", pdbBreakRegex, pdbSteptoLocRegex)
	case "lldb":
		return newSteptoAdapter(editorState, "lldb", "(lldb) ", lldbBreakRegex, lldbSteptoLocRegex)
	case "dlv":
		return newSteptoAdapter(editorState, "dlv", "(dlv) ", dlvBreakRegex, dlvSteptoLocRegex)
	}
	log.Fatalln("Invalid debugger tool")
	return nil
}
