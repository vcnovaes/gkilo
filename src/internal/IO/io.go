package ioc

import (
	"bufio"
	"os"

	textbuffer "github.com/vcnovaes/gkilo/src/internal/core/TextBuffer"

	"github.com/nsf/termbox-go"
)

func LoadInputFile(defaultfile string) (*os.File, string, error) {
	sourcefile := defaultfile
	if len(os.Args) > 1 {
		sourcefile = os.Args[1]
	}
	file, err := os.Open(sourcefile)
	if err != nil {
		file, err = os.Create(sourcefile)
	}
	return file, sourcefile, err
}

func WriteFile(sourcefile string, buffer textbuffer.TextBuffer) error {
	file, err := os.Create(sourcefile)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	for _, line := range buffer.GetBuffer() {
		linetoWrite := string(line) + "\n"
		writer.WriteString(linetoWrite)
	}
	writer.Flush()
	defer file.Close()
	return nil
}

func GetPressedKey() termbox.Event {
	var keyEvent termbox.Event
	switch event := termbox.PollEvent(); event.Type {
	case termbox.EventKey:
		keyEvent = event
	case termbox.EventError:
		panic(event.Err)
	}
	return keyEvent
}
