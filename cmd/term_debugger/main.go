package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/fatih/color"
	"github.com/jiangyinzuo/term-debugger/package/debugger"
	"github.com/jiangyinzuo/term-debugger/package/texteditor"
)

type debugWrapper struct {
	wg         sync.WaitGroup
	debugger   debugger.TermDebugger
	textEditor texteditor.TextEditor
}

func startDebugWrapperAndWait(cmd *exec.Cmd) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	wrapper := &debugWrapper{
		debugger:   &debugger.GDBAdapter{},
		textEditor: &texteditor.EscapeSequence{},
	}
	wrapper.wg.Add(2)
	go wrapper.processChildStderr(stderr)
	go wrapper.processChildStdout(stdout)
	go wrapper.processChildStdin(stdin)

	cmd.Run()
	wrapper.wg.Wait()
	return
}

func (d *debugWrapper) processChildStdin(stdin io.WriteCloser) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		d.debugger.ProcessUserInput(stdin, scanner.Text())
	}
}

func (d *debugWrapper) processChildStderr(stderr io.Reader) {
	defer d.wg.Done()
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		d.debugger.ProcessChildStderr(scanner.Text(), d.textEditor)
	}
}

func (d *debugWrapper) processChildStdout(stdout io.Reader) {
	defer d.wg.Done()

	buffer := make([]byte, 1)
	var output bytes.Buffer

	suffix := d.debugger.EndPromt()
	for {
		_, err := stdout.Read(buffer)
		if err != nil {
			break
		}
		char := buffer[0]
		output.WriteByte(char)

		if output.Len() > len(suffix) && bytes.HasSuffix(output.Bytes(), []byte(suffix)) {
			d.debugger.ProcessChildStdout(output.String(), d.textEditor)
			output.Reset()
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: " + os.Args[0] + " <command>")
		os.Exit(1)
	}
	blue := color.New(color.FgBlue)
	blue.Println("Running command: ", os.Args[1:])
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	startDebugWrapperAndWait(cmd)
	blue.Println("Bye " + os.Args[1] + "!")
}

func handleUserInput(stdin io.WriteCloser, handler debugger.TermDebugger) {
}
