package debugger

import (
	"fmt"
	"io"
	"strings"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
)

type GDBAdapter struct {
}

func (h *GDBAdapter) Name() string {
	return "gdb"
}
func (h *GDBAdapter) EndPromt() string {
	return "(gdb) "
}

func (h *GDBAdapter) ProcessChildStdout(output string, editor texteditor.TextEditor) {
	fmt.Println("[STDOUT]")
	fmt.Print(output)
	editor.SendKey("")
}

func (h *GDBAdapter) ProcessUserInput(stdin io.WriteCloser, input string) {
	fmt.Println("You entered:", input)
	stdin.Write([]byte(input + "\n"))
}

func (h *GDBAdapter) ProcessChildStderr(output string, editor texteditor.TextEditor) {
	fmt.Println("[STDERR]")
	fmt.Println(output)
}

func contains(str, substr string) bool {
	return strings.Contains(str, substr)
}
