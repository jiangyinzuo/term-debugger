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
	log "github.com/sirupsen/logrus"
)

type debugWrapper struct {
	wg       sync.WaitGroup
	debugger debugger.TermDebugger
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
	textEditor := &texteditor.VimTapi{}
	editorState := texteditor.NewEditorState(textEditor)
	cmd.Stderr = cmd.Stdout
	wrapper := &debugWrapper{
		debugger: debugger.NewGDBAdapter(editorState),
	}

	wrapper.wg.Add(1)
	go wrapper.processChildOutput(stdout)
	go wrapper.processChildStdin(stdin)

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	wrapper.wg.Wait()
	return
}

func (d *debugWrapper) processChildStdin(stdin io.WriteCloser) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		d.debugger.ProcessUserInput(stdin, scanner.Text())
	}
}

func (d *debugWrapper) processChildOutput(stdout io.Reader) {
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

		if char == '\n' || (output.Len() >= len(suffix) && bytes.HasSuffix(output.Bytes(), []byte(suffix))) {
			d.debugger.ProcessChildOutput(output.String())
			output.Reset()
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: " + os.Args[0] + " <command>")
		os.Exit(1)
	}
	log.SetLevel(log.DebugLevel)
	blue := color.New(color.FgBlue)
	blue.Println("Running command: ", os.Args[1:])
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	startDebugWrapperAndWait(cmd)
	blue.Println("Bye " + os.Args[1] + "!")
}
