package texteditor

import "fmt"

type EscapeSequence struct {
}

func (e *EscapeSequence) SendKey(key string) {
	fmt.Printf("\x1b]51;[\"drop\", \"README.md\"]\x07")
}
