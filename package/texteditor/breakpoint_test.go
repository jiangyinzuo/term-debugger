package texteditor_test

import (
	"testing"

	"github.com/jiangyinzuo/term-debugger/package/texteditor"
	"github.com/stretchr/testify/assert"
)

func TestBreakpointMap(t *testing.T) {
	bpm := texteditor.NewBreakpointMap()
	bpm.Add(1, "main.go", 10)
	bpm.Add(1, "main.go", 10)
	bpm.Add(2, "main.go", 10)

	filename, line, res := bpm.GetByID(1)
	assert.True(t, res)
	assert.Equal(t, "main.go", filename)
	assert.Equal(t, 10, line)
	bpm.RemoveAll()
	filename, line, res = bpm.GetByID(1)
	assert.False(t, res)
}
