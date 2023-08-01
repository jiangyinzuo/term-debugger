package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/fatih/color"
	"github.com/jiangyinzuo/term-debugger/package/debugger"
	log "github.com/sirupsen/logrus"
)

var debugTool = flag.String("debug", "gdb", "debug tool to use")
var logLevel = flag.String("log", "info", "log level")

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
	cmd.Stderr = cmd.Stdout
	wrapper := &debugWrapper{
		debugger: debugger.NewDebugger(debugTool),
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
	flag.Parse()
	logLevel, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	blue := color.New(color.FgBlue)
	args := flag.Args()
	if len(args) == 0 {
		blue.Println("command not found")
		os.Exit(1)
	}
	blue.Println("Running command: ", args)
	cmd := exec.Command(args[0], args[1:]...)

	startDebugWrapperAndWait(cmd)
	blue.Println("Bye " + args[0] + "!")
}
