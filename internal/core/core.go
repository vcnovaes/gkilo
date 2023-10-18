package core

import (
	"fmt"
	"ioc"
	"os"
	"render"
	"textbuffer"

	"github.com/nsf/termbox-go"
)

type Mode uint8

const (
	MODE_DEFAULT Mode = iota
	MODE_EDIT
)

type Editor struct {
	buffer textbuffer.TextBuffer
	canvas render.Canvas
	file   File
	mode   Mode
}
type File struct {
	filename string
	modified bool
}

func (e *Editor) Init() {
	err := termbox.Init()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defaultfilename := "unamed"
	file, filename, err := ioc.LoadInputFile(defaultfilename)
	e.file.filename = filename
	if err != nil {
		fmt.Println("Error on inputfile:", err)
		os.Exit(1)
	}
	e.buffer.LoadFile(file)
}

func (e Editor) Close() {
	termbox.Close()
	os.Exit(0)
}

func (e *Editor) Run() {
	theme := render.Theme{
		LineNumber:       render.Composition{FG: termbox.ColorLightMagenta, BG: termbox.ColorDefault},
		Text:             render.Composition{FG: termbox.ColorGreen, BG: termbox.ColorDefault},
		EmptyLine:        render.Composition{FG: termbox.ColorDefault, BG: termbox.ColorDefault},
		EditorBackground: termbox.ColorDefault,
		EmptyLineSymbol:  "~",
	}

	for {
		e.canvas.Init(theme)
		e.buffer.Update(&e.canvas)
		e.canvas.Flush()
		e.processKeyEvent()
	}
}

func (e *Editor) writeFile() {
	ioc.WriteFile(e.file.filename, e.buffer)
}

func (e *Editor) handleCommand(ch rune) {
	switch ch {
	case 'q':
		e.Close()
	case 'e':
		e.mode = MODE_EDIT
	case 'w':
		e.writeFile()
	case 'j':
		e.buffer.MoveCursorDown()
	case 'k':
		e.buffer.MoveCursorUp()
	case 'l':
		e.buffer.MoveCursorLeft()
	case 'h':
		e.buffer.MoveCursorRight()
	}
}

func (e *Editor) handleNoCharEditInput(key termbox.Event) {
	switch key.Key {
	case termbox.KeySpace:
		e.buffer.Write(key.Ch)
	case termbox.KeyEnter:
		e.buffer.BreakLine()
	case termbox.KeyCtrlQ:
		e.Close()
	case termbox.KeyBackspace:
	case termbox.KeyBackspace2:
		e.buffer.Erase()
	default:
		break
	}
}
func (e *Editor) handleNoChar(key termbox.Key) {
	switch key {
	case termbox.KeyArrowRight:
		e.buffer.MoveCursorRight()
	case termbox.KeyArrowLeft:
		e.buffer.MoveCursorLeft()
	case termbox.KeyArrowDown:
		e.buffer.MoveCursorDown()
	case termbox.KeyArrowUp:
		e.buffer.MoveCursorUp()
	}
}

func (e *Editor) processKeyEvent() {
	keyEvent := ioc.GetPressedKey()
	if keyEvent.Key == termbox.KeyEsc {
		e.mode = MODE_DEFAULT
	} else if keyEvent.Ch != 0 {
		if e.mode == MODE_EDIT {
			e.buffer.Write(keyEvent.Ch)
			e.file.modified = true
			return
		}
		e.handleCommand(keyEvent.Ch)
	} else {
		if e.mode == MODE_EDIT {
			e.handleNoCharEditInput(keyEvent)
			return
		}
		e.handleNoChar(keyEvent.Key)

		e.buffer.UpdateCol()
	}
}
