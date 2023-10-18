package textbuffer

import (
	"bufio"
	"os"
	"render"
)

type TextBuffer struct {
	buffer  []BufferLine
	currRow int
	currCol int

	offsetRow int
	offsetCol int
}

type BufferLine []rune

func (tb TextBuffer) GetBuffer() []BufferLine {
	return tb.buffer
}
func (tb TextBuffer) currentLine() []rune {
	return tb.buffer[tb.currRow]
}

func (tb *TextBuffer) UpdateCol() {
	if tb.currCol > len(tb.currentLine()) {
		tb.currCol = len(tb.currentLine())
	}
}

func (tb *TextBuffer) Write(char rune) {
	bufferLine := make(BufferLine, len(tb.currentLine())+1)
	copy(bufferLine[:tb.currCol], tb.currentLine()[:tb.currCol])
	bufferLine[tb.currCol] = char
	copy(bufferLine[tb.currCol+1:], tb.currentLine()[tb.currCol:])
	tb.buffer[tb.currRow] = bufferLine
	tb.currCol++
}

func (tb *TextBuffer) BreakLine() {

	firstHalf := make([]rune, len(tb.currentLine()[:tb.currCol]))
	secondHalf := make([]rune, len(tb.currentLine()[tb.currCol:]))

	copy(firstHalf, tb.buffer[tb.currRow][:tb.currCol])
	copy(secondHalf, tb.buffer[tb.currRow][tb.currCol:])

	newBuffer := make([]BufferLine, len(tb.buffer)+1)
	newBuffer[tb.currRow] = firstHalf
	tb.currRow++
	newBuffer[tb.currRow] = secondHalf
	tb.currCol = 0

	copy(newBuffer[:tb.currRow], tb.buffer[:tb.currRow-1])
	copy(newBuffer[tb.currRow:], tb.buffer[tb.currRow:])

	tb.buffer = newBuffer
}

func (tb *TextBuffer) eraseCol() {
	tb.currCol--
	alteredRow := make([]rune, len(tb.currentLine())-1)
	copy(alteredRow[:tb.currCol], tb.currentLine()[:tb.currCol])
	copy(alteredRow[tb.currCol:], tb.currentLine()[tb.currCol+1:])
	tb.buffer[tb.currRow] = alteredRow
}

func (tb *TextBuffer) eraseLine() {
	appendedLine := make([]rune, len(tb.currentLine()))
	copy(appendedLine, tb.currentLine()[tb.currCol:])

	newTextBuffer := make([]BufferLine, len(tb.buffer)-1)

	copy(newTextBuffer[:tb.currRow], tb.buffer[:tb.currRow])
	copy(newTextBuffer[tb.currRow:], tb.buffer[tb.currRow+1:])

	tb.buffer = newTextBuffer
	tb.currRow--

	tb.currCol = len(tb.currentLine())

	updatedLine := make([]rune, len(tb.buffer)+len(appendedLine))

	copy(updatedLine[:len(tb.currentLine())], tb.currentLine())
	copy(updatedLine[len(tb.currentLine()):], appendedLine)

	tb.buffer[tb.currRow] = updatedLine
}

func (tb *TextBuffer) Erase() {
	if tb.currCol > 0 {
		tb.eraseCol()
	} else if tb.currRow > 0 {
		tb.eraseLine()
	}
}

func (tb *TextBuffer) Display(canvas *render.Canvas) {
	var col int
	for row := 0; row < canvas.Height; row++ {
		tbRow := row + tb.offsetRow
		for col = 0; col < canvas.Width; col++ {
			tbCol := col + tb.offsetCol
			if tbRow < len(tb.buffer) {
				canvas.RenderLineNumber(row, tb.offsetRow)
				if tbRow >= 0 && tbCol < len(tb.buffer[tbRow]) {
					symbol := tb.buffer[tbRow][tbCol]
					if symbol == rune('\t') {
						symbol = rune(' ')
					}
					canvas.RenderRune(col, row, symbol)
				}
			} else if row+tb.offsetRow > (len(tb.buffer) - 1) {
				canvas.RenderEmptyRow(row)
			}
		}
		canvas.RenderLineBreak(col, row)
	}
}

func (tb *TextBuffer) setCursorPosition(canvas *render.Canvas) {
	canvas.RenderCursor(tb.currCol-tb.offsetCol, tb.currRow-tb.offsetRow)
}

func (tb *TextBuffer) EnableScroll(canvas *render.Canvas) {
	if tb.currRow < tb.offsetRow {
		tb.offsetRow = tb.currRow
	}
	if tb.currRow >= tb.offsetRow+canvas.Height {
		tb.offsetRow = tb.currRow - canvas.Height + 1
	}

	if tb.currCol < tb.offsetCol {
		tb.offsetCol = tb.currCol
	}
	if tb.currCol >= tb.offsetCol+canvas.Width-canvas.LineNumberWidth {
		tb.offsetCol = tb.currCol - canvas.Width + canvas.LineNumberWidth + 1
	}
}

func (tb *TextBuffer) Update(canvas *render.Canvas) {
	tb.EnableScroll(canvas)
	tb.Display(canvas)
	tb.setCursorPosition(canvas)
}

func (tb *TextBuffer) LoadFile(file *os.File) {
	lineIdx := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scannedLine := scanner.Text()
		tb.buffer = append(tb.buffer, BufferLine{})
		for _, ch := range scannedLine {
			tb.buffer[lineIdx] = append(tb.buffer[lineIdx], rune(ch))
		}
		lineIdx++
	}
	if lineIdx <= 1 {
		tb.buffer = append(tb.buffer, BufferLine{})
	}
}
func (tb *TextBuffer) SetBuffer(buffer []BufferLine) {
	tb.buffer = buffer
}

func (tb *TextBuffer) MoveCursorUp() {
	if tb.currRow != 0 {
		tb.currRow--
	}
}

func (tb *TextBuffer) MoveCursorDown() {
	if tb.currRow < (len(tb.buffer) - 1) {
		tb.currRow++
	}
}

func (tb *TextBuffer) MoveCursorLeft() {
	if tb.currCol > 0 {
		tb.currCol--
	}
}

func (tb *TextBuffer) MoveCursorRight() {
	if tb.currCol < (len(tb.currentLine()) - 1) {
		tb.currCol++
	}
}
