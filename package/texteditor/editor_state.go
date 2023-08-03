package texteditor

type EditorState struct {
	breakpointMap *BreakpointMap

	currentFileName string
	currentLine     int

	textEditor TextEditor
}

func NewEditorState(textEditor TextEditor) *EditorState {
	return &EditorState{
		breakpointMap:   NewBreakpointMap(),
		currentFileName: "",
		currentLine:     0,
		textEditor:      textEditor,
	}
}

func (es *EditorState) SetCurrentFileName(fileName string) {
	es.currentFileName = fileName
}

func (es *EditorState) SetCurrentLine(line int) {
	es.currentLine = line
}

func (es *EditorState) StepToCursor() {
	es.textEditor.StepTo(es.currentFileName, es.currentLine)
}

func (es *EditorState) AddBreakpoint(id int, fileName string, line int) {
	es.breakpointMap.Add(id, fileName, line)
	es.textEditor.SetBreakpoint(es.breakpointMap.CalculateBrkType(fileName, line), fileName, line)
}

func (es *EditorState) RemoveBreakpointByLoc(fileName string, line int) {
	es.breakpointMap.RemoveByLocation(fileName, line)
}

func (es *EditorState) RemoveBreakpointByLine(line int) {
	es.breakpointMap.RemoveByLocation(es.currentFileName, line)
}

func (es *EditorState) RemoveBreakpointById(id int) {
	es.breakpointMap.RemoveByID(id)
}

func (es *EditorState) RemoveAllBreakPoints() {
	es.breakpointMap.RemoveAll()
	es.textEditor.
}
