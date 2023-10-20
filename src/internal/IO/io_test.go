package ioc_test

import (
	"os"
	"testing"

	ioc "github.com/vcnovaes/gkilo/src/internal/IO"
	textbuffer "github.com/vcnovaes/gkilo/src/internal/core/TextBuffer"
)

func TestLoadInputFile(t *testing.T) {
	t.Skip("Break test")
	// Create a temporary file for testing
	tempFile, err := os.Create("testfile")
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the temporary file

	// Test loading the existing file
	file, _, err := ioc.LoadInputFile(tempFile.Name())
	if err != nil {
		t.Errorf("LoadInputFile() failed with an existing file: %v", err)
	}
	if file == nil {
		t.Errorf("LoadInputFile() returned a nil file")
	}
	file.Close()

	// Test creating a new file
	newFile, _, newErr := ioc.LoadInputFile("nonexistent.txt")
	if newErr != nil {
		t.Errorf("LoadInputFile() failed when creating a new file: %v", newErr)
	}
	if newFile == nil {
		t.Errorf("LoadInputFile() returned a nil file when creating a new file")
	}
	newFile.Close()
}

func TestWriteFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	buffer := textbuffer.TextBuffer{}
	var testlines [][]rune
	testlines = make([][]rune, 2)
	testStrs := []string{"Line 1", "Line 2"}
	for i, s := range testStrs {
		for _, ch := range s {
			testlines[i] = append(testlines[i], ch)
		}
	}
	buffer.SetBuffer([]textbuffer.BufferLine{testlines[0], testlines[1]})

	err = ioc.WriteFile(tempFile.Name(), buffer)
	if err != nil {
		t.Errorf("WriteFile() failed: %v", err)
	}

	fileContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read the temporary file: %v", err)
	}

	expectedContent := []byte("Line 1\nLine 2\n")
	if string(fileContent) != string(expectedContent) {
		print(fileContent)
		t.Errorf("WriteFile() wrote incorrect content, expected: %s, got: %s", string(expectedContent), string(fileContent))
	}
}
