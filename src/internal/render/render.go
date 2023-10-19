package render

import (
	"strconv"

	"github.com/vcnovaes/gkilo/src/internal/editorconfig"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Theme struct {
	LineNumber       Composition
	Text             Composition
	EmptyLine        Composition
	EditorBackground termbox.Attribute
	EmptyLineSymbol  string
}

func ColorMap(color string) termbox.Attribute {
	colormap := make(map[string]termbox.Attribute)
	colormap["blue"] = (termbox.ColorBlue)
	colormap["lightblue"] = (termbox.ColorLightBlue)
	colormap["green"] = (termbox.ColorGreen)
	colormap["lightgreen"] = (termbox.ColorLightBlue)
	colormap["black"] = (termbox.ColorBlack)
	colormap["default"] = (termbox.ColorDefault)
	colormap["magenta"] = (termbox.ColorMagenta)
	colormap["lightmagenta"] = (termbox.ColorLightMagenta)
	colormap["red"] = (termbox.ColorRed)
	colormap["lightred"] = (termbox.ColorLightRed)
	colormap["white"] = (termbox.ColorWhite)

	return colormap[color]

}

type Composition struct {
	FG, BG termbox.Attribute
}

func (c *Composition) LoadFromString(fg, bg string) {
	c.BG = ColorMap(bg)
	c.FG = ColorMap(fg)
}

func (c *Composition) LoadFromStringComposition(cstr editorconfig.CompositionStr) {
	c.BG = ColorMap(cstr.BG)
	c.FG = ColorMap(cstr.FG)
}

type Canvas struct {
	theme           Theme
	Width, Height   int
	LineNumberWidth int
}

func render(x, y int, fg, bg termbox.Attribute, content string) {
	for _, char := range content {
		termbox.SetCell(x, y, char, fg, bg)
		x += runewidth.RuneWidth(char)
	}
}

func (c *Canvas) RenderLineNumber(row int, offsetRow int) {
	digit := strconv.Itoa(row + offsetRow + 1)
	offset := c.LineNumberWidth - len(digit) - 1
	render(offset, row, c.theme.LineNumber.FG, c.theme.LineNumber.BG, digit)
}

func (c *Canvas) RenderEmptyRow(row int) {
	render(0, row, c.theme.EmptyLine.FG, c.theme.EmptyLine.BG,
		c.theme.EmptyLineSymbol)
}

func (c *Canvas) RenderRune(col, row int, char rune) {
	render(
		col+c.LineNumberWidth,
		row,
		c.theme.Text.FG, c.theme.LineNumber.BG,
		string(char))
}

func (c *Canvas) RenderLineBreak(col, row int) {
	render(
		col,
		row,
		c.theme.Text.FG, c.theme.Text.BG, "\n")
}

func (c Canvas) RenderCursor(col, row int) {
	termbox.SetCursor(col+c.LineNumberWidth, row)
}

func (c Canvas) Clear() {
	termbox.Clear(c.theme.EditorBackground, c.theme.EditorBackground)
}

func (c Canvas) Flush() { termbox.Flush() }

func (c *Canvas) Init(theme Theme) {
	c.Width, c.Height = termbox.Size()
	c.theme = theme
	c.LineNumberWidth = 4
	c.Clear()
}
