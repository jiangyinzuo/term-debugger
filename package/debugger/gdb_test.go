package debugger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegex(t *testing.T) {
	testcase := []struct {
		regex    string
		filename string
		line     string
	}{
		{"Breakpoint 1, main () at hello.cpp:22\n", "hello.cpp", "22"},
		{"foo (a=2) at hello.cpp:16\n", "hello.cpp", "16"},
	}
	for _, tc := range testcase {
		matches := stepToCursorLocRegex.FindStringSubmatch(tc.regex)
		assert.Equal(t, 3, len(matches))
		assert.Equal(t, tc.filename, matches[1])
		assert.Equal(t, tc.line, matches[2])
	}
}
