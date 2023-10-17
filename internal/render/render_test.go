package render

import (
	"testing"

	"github.com/nsf/termbox-go"
)

func TestCanvasInit(t *testing.T) {
	c := &Canvas{}
	theme := Theme{
		LineNumber: Composition{
			FG: termbox.ColorWhite,
			BG: termbox.ColorBlack,
		},
	}

	c.Init(theme)

	if c.LineNumberWidth != 4 {
		t.Errorf("Canvas Init did not set LineNumberWidth correctly")
	}
	if c.theme != theme {
		t.Errorf("Canvas Init did not set the theme correctly")
	}
}

// You can add more test cases as needed to cover other functions in the package.
