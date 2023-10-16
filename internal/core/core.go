package core

import (
	"bufio"
	"common"
	"fmt"
	tinterface "interface"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type MODE uint8

const (
	MODE_DEFAULT = iota
	MODE_EDIT
)

type EditorCtx struct {
	originTermios        *syscall.Termios
	screenRows           int
	screenCols           int
	sourceFile           string
	currentRow           int
	currentCol           int
	textBuffer           [][]rune
	lineNumberWidth      int
	offsetCol, offsetRow int
	modifiedFile         bool
	mode                 MODE
}
type Editor struct {
	ctx EditorCtx
}

func (e Editor) getCurrentRow() []rune {
	return e.ctx.textBuffer[e.ctx.currentRow]
}

func (e *Editor) enableScrollTextBuffer() {
	if e.ctx.currentRow < e.ctx.offsetRow {
		e.ctx.offsetRow = e.ctx.currentRow
	}
	if e.ctx.currentRow >= e.ctx.offsetRow+e.ctx.screenRows {
		e.ctx.offsetRow = e.ctx.currentRow - e.ctx.screenRows + 1
	}

	if e.ctx.currentCol < e.ctx.offsetCol {
		e.ctx.offsetCol = e.ctx.currentCol
	}
	if e.ctx.currentCol >= e.ctx.offsetCol+e.ctx.screenCols-e.ctx.lineNumberWidth {
		e.ctx.offsetCol = e.ctx.currentCol - e.ctx.screenCols + e.ctx.lineNumberWidth + 1
	}
}
func (e *Editor) defaultFile() {
	e.ctx.sourceFile = "unamed"
	e.ctx.textBuffer = append(e.ctx.textBuffer, []rune{})
}

func (e *Editor) loadFile() {
	if len(os.Args) > 1 {
		sourceFile := os.Args[1]
		ReadFile(sourceFile, &e.ctx)
	} else {
		e.defaultFile()
	}
}
func (e Editor) getBufferSize() int {
	return len(e.ctx.textBuffer)
}

func (e Editor) render(x, y int, foreground, backaground termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(x, y, ch, foreground, backaground)
		x += runewidth.RuneWidth(ch)
	}
}

func (e *Editor) getLineNumberWidth() int {
	e.ctx.lineNumberWidth = len(strconv.Itoa(e.getBufferSize())) + 1
	return e.ctx.lineNumberWidth
}

func (e *Editor) displayTextBuffer() {
	for row := 0; row < e.ctx.screenRows; row++ {
		textBufferRow := row + e.ctx.offsetRow
		for col := 0; col < e.ctx.screenCols; col++ {
			textBufferCol := col + e.ctx.offsetCol
			if textBufferRow < e.getBufferSize() {
				lineNumberOffset := e.getLineNumberWidth() - len(strconv.Itoa(textBufferRow+1))
				e.render(lineNumberOffset, row, termbox.ColorWhite, termbox.ColorDarkGray,
					strconv.Itoa(textBufferRow+1))
			}
			if textBufferRow >= 0 && textBufferRow < e.getBufferSize() &&
				textBufferCol < len(e.ctx.textBuffer[textBufferRow]) {
				if e.ctx.textBuffer[textBufferRow][textBufferCol] != rune('\t') {
					termbox.SetCell(col+e.ctx.lineNumberWidth+1, row,
						e.ctx.textBuffer[textBufferRow][textBufferCol],
						termbox.ColorDefault, termbox.ColorDefault)
				} else {
					termbox.SetCell(col+e.ctx.lineNumberWidth, row, rune(' '), termbox.ColorGreen, termbox.ColorDefault)
				}
			} else if row+e.ctx.offsetRow > (e.getBufferSize() - 1) {
				termbox.SetCell(0, row, '~', termbox.ColorDarkGray, termbox.ColorDefault)
			}
		}
		termbox.SetChar(e.ctx.screenCols-1, row, '\n')
	}
}

func (e *Editor) setCursorPosition() {
	termbox.SetCursor(e.ctx.currentCol-e.ctx.offsetCol+e.ctx.lineNumberWidth,
		e.ctx.currentRow-e.ctx.offsetRow)
}

func (e *Editor) setCurrentRow(insertRow []rune) {
	e.ctx.textBuffer[e.ctx.currentRow] = insertRow
	e.ctx.currentCol++
}

func (e *Editor) Init() {
	err := termbox.Init()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	e.loadFile()
}

func (e Editor) Close() {
	termbox.Close()
}

func (e *Editor) insertRune(event termbox.Event) {
	insertLine := make([]rune, len(e.ctx.textBuffer[e.ctx.currentRow])+1)
	copy(insertLine[:e.ctx.currentCol], e.ctx.textBuffer[e.ctx.currentRow][:e.ctx.currentCol])
	ch := rune(event.Ch)
	if event.Key == termbox.KeySpace || event.Key == termbox.KeyTab {
		ch = rune(' ')
	}
	insertLine[e.ctx.currentCol] = ch
	copy(insertLine[e.ctx.currentCol+1:], e.ctx.textBuffer[e.ctx.currentRow][e.ctx.currentCol:])
	e.ctx.textBuffer[e.ctx.currentRow] = insertLine
	e.ctx.currentCol++
}

func (e *Editor) Run() {
	for {

		e.ctx.screenCols, e.ctx.screenRows = termbox.Size()
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		e.enableScrollTextBuffer()
		e.displayTextBuffer()
		e.setCursorPosition()
		termbox.Flush()
		e.processKeyEvent()
	}
}

func (e *Editor) addLine() {

	// spliting in two halfs
	firstHalf := make([]rune, len(e.getCurrentRow()[e.ctx.currentCol:]))
	secondHalf := make([]rune, len(e.getCurrentRow()[:e.ctx.currentCol]))
	// copying the content
	copy(firstHalf, e.ctx.textBuffer[e.ctx.currentRow][:e.ctx.currentCol])
	copy(secondHalf, e.ctx.textBuffer[e.ctx.currentRow][e.ctx.currentCol:])

	e.ctx.textBuffer[e.ctx.currentRow] = firstHalf

	e.ctx.currentRow++
	e.ctx.currentCol = 0

	newBuffer := make([][]rune, len(e.ctx.textBuffer)+1)
	copy(newBuffer[:e.ctx.currentRow], e.ctx.textBuffer[:e.ctx.currentRow])
	newBuffer[e.ctx.currentRow] = secondHalf

	copy(newBuffer[e.ctx.currentRow+1:], e.ctx.textBuffer[e.ctx.currentRow:])

	e.ctx.textBuffer = newBuffer
}

func (e *Editor) processKeyEvent() {
	keyEvent := getKey()
	if keyEvent.Key == termbox.KeyEsc {
		e.ctx.mode = MODE_DEFAULT
	} else if keyEvent.Ch != 0 {
		if e.ctx.mode == MODE_EDIT {
			e.insertRune(keyEvent)
			e.ctx.modifiedFile = true
			if keyEvent.Ch == 'q' {
				defer func() {
					e.Close()
					os.Exit(0)
				}()
			}
			return
		}
		switch keyEvent.Ch {
		case 'q':
			termbox.Close()
			os.Exit(0)
		case 'e':
			print("EDIT")
			e.ctx.mode = MODE_EDIT
		}
	} else {
		switch keyEvent.Key {
		case termbox.KeySpace:
			e.insertRune(keyEvent)
		case termbox.KeyEnter:
			e.addLine()
		case termbox.KeyCtrlQ:
			e.Close()
			os.Exit(0)
		}

	}
}

func getKey() termbox.Event {
	var event termbox.Event
	switch event = termbox.PollEvent(); event.Type {
	case termbox.EventKey:
		return event
	case termbox.EventError:
		panic(event.Err)
	}
	return event
}

const CTRL_Q byte = 17

func getTermios(fd uintptr) *syscall.Termios {
	var t syscall.Termios
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL, // Input output control
		os.Stdin.Fd(),
		syscall.TCGETS,
		uintptr(unsafe.Pointer(&t)),
		0, 0, 0)

	if err != 0 {
		panic("Error getting termios")
	}

	return &t
}

func setTermios(fd uintptr, term *syscall.Termios) {
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		os.Stdin.Fd(),
		syscall.TCSETS,
		uintptr(unsafe.Pointer(term)),
		0, 0, 0)
	if err != 0 {
		panic("err")
	}
}

func setRaw(term *syscall.Termios) {
	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	term.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	// newState.Oflag &^= syscall.OPOST
	term.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	term.Cflag &^= syscall.CSIZE | syscall.PARENB
	term.Cflag |= syscall.CS8

	term.Cc[syscall.VMIN] = 1
	term.Cc[syscall.VTIME] = 0
}

func setupRawMode(editorConfig *EditorCtx) {
	editorConfig.originTermios = getTermios(os.Stdin.Fd())
	defer setTermios(
		os.Stdin.Fd(),
		editorConfig.originTermios)

	setRaw(editorConfig.originTermios)
	setTermios(os.Stdin.Fd(), editorConfig.originTermios)
}

func runEditor() {
	for {
		tinterface.EditorRefreshScreen()
		processKeyPressed()
	}
}

func processKeyPressed() {
	buffer := make([]byte, 1)
	syscall.Read(0, buffer)
	ch := string(buffer)
	print(ch)

	if buffer[0] == CTRL_Q {
		tinterface.EditorRefreshScreen()
		os.Exit(0)
	}
}

func initEditor(editorConfig *EditorCtx) {
	editorConfig.screenCols, editorConfig.screenRows, _ = common.GetWindowSize()
	setupRawMode(editorConfig)
}

func ReadFile(filename string, editorCtx *EditorCtx) {
	file, err := os.Open(filename)
	editorCtx.sourceFile = filename
	if err != nil {
		editorCtx.textBuffer = append(editorCtx.textBuffer, []rune{})
		return
	}
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		scannedLine := scanner.Text()
		editorCtx.textBuffer = append(editorCtx.textBuffer, []rune{})
		for _, ch := range scannedLine {
			editorCtx.textBuffer[lineNumber] = append(editorCtx.textBuffer[lineNumber], rune(ch))
		}
		lineNumber++
	}
	if lineNumber == 0 {
		editorCtx.textBuffer = append(editorCtx.textBuffer, []rune{})
	}
}
