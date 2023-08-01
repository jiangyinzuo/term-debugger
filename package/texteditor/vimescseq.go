package texteditor

import (
	"encoding/json"
	"fmt"
)

type VimTapi struct {
}

func (v *VimTapi) AddBreakPoint(id int, filename string, line int) {
	v.callTapi("Tapi_TermDebuggerAddBreakpoint", []interface{}{filename, line})
}

func (v *VimTapi) DeleteBreakPointByLoc(filename string, line int) {
	v.callTapi("Tapi_TermDebuggerDeleteBreakPointByLoc", []interface{}{filename, line})
}

func (v *VimTapi) DeleteBreakPointByID(id int) {
	v.callTapi("Tapi_TermDebuggerDeleteBreakPointByID", []interface{}{id})
}

func (v *VimTapi) DeleteAllBreakPoints() {
	v.callTapi("Tapi_TermDebuggerDeleteAllBreakPoints", []interface{}{})
}

func (v *VimTapi) callTapi(funcName string, args []interface{}) {
	b, err := json.Marshal([]interface{}{"call", funcName, args})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\x1b]51;%s\x07", string(b))
}

