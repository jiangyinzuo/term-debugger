package texteditor

import (
	"encoding/json"
	"fmt"
)

type VimTapi struct {
}

func (v *VimTapi) StepTo(filename string, line int) {
	v.callTapi("Tapi_TermDebuggerStepTo", []interface{}{filename, line})
}

func (v *VimTapi) SignBreakpoint(brkType byte, filename string, line int) {
	v.callTapi("Tapi_TermDebuggerSignBreakpoint", []interface{}{brkType, filename, line})
}

func (v *VimTapi) callTapi(funcName string, args []interface{}) {
	b, err := json.Marshal([]interface{}{"call", funcName, args})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\x1b]51;%s\x07", string(b))
}
