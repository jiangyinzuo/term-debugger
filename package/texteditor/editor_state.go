package texteditor

type EditorState struct {
	currentFileName string
	currentLine     int

	textEditor TextEditor
}

func NewEditorState(textEditor TextEditor) *EditorState {
	return &EditorState{
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

func (es *EditorState) SyncCursor() {
	if es.currentFileName == "" {
		return
	}
	es.textEditor.StepTo(es.currentFileName, es.currentLine)
}

func (es *EditorState) SignBreakpoint(brkType byte, filename string, line int) {
	es.textEditor.SignBreakpoint(brkType, filename, line)
}

