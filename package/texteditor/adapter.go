package texteditor

type TextEditor interface {
	AddBreakPoint(id int, filename string, line int)
	DeleteBreakPointByLoc(filename string, line int)
	DeleteBreakPointByID(id int)
	DeleteAllBreakPoints()
}
