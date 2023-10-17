package textbuffer

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	tb := TextBuffer{}
	tb.SetBuffer([]BufferLine{{}})
	tb.Write('a')
	content := string(tb.currentLine())
	print(content)
	if string(tb.currentLine()) != "a" {
		t.Errorf("Write() did not update the current line correctly")
	}
}

func TestErase(t *testing.T) {
	tb := TextBuffer{}
	tb.SetBuffer([]BufferLine{{}})
	tb.Write('a')
	tb.Erase()
	content := string(tb.currentLine())
	print(content)

	if len(tb.currentLine()) != 0 || tb.currCol != 0 {
		t.Errorf("Erase() did not clear the buffer correctly")
	}
}

func TestLoadFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the temporary file

	// Write some content to the temporary file
	content := "Hello, World!\nThis is a test.\n"
	tempFile.WriteString(content)
	tempFile.Close()

	tb := &TextBuffer{}
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to open the temporary file: %v", err)
	}
	tb.LoadFile(file)

	expectedBuffer := []BufferLine{
		[]rune("Hello, World!"),
		[]rune("This is a test."),
	}

	if len(tb.buffer) != len(expectedBuffer) {
		t.Errorf("LoadFile() did not load the correct number of lines")
	} else {
		for i, line := range tb.buffer {
			if string(line) != string(expectedBuffer[i]) {
				t.Errorf("LoadFile() loaded incorrect content in line %d", i)
			}
		}
	}
}

// You can add more test cases for other methods as needed.
