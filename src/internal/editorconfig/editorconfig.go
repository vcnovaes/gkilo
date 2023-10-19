package editorconfig

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

const CONFIG_FILE = "config.json"

type Config struct {
	Theme Theme `json:"Theme"`
}

type Theme struct {
	LineNumber       CompositionStr `json:"LineNumber"`
	Text             CompositionStr `json:"Text"`
	EmptyLine        CompositionStr `json:"EmptyLine"`
	EditorBackground string         `json:"EditorBackground"`
	EmptyLineSymbol  string         `json:"EmptyLineSymbol"`
}

type CompositionStr struct {
	FG string `json:"fg"`
	BG string `json:"bg"`
}

func GetConfig() Config {
	var config Config
	jsonfile, _ := os.Open(CONFIG_FILE)

	byteValue, _ := io.ReadAll(jsonfile)
	json.Unmarshal(byteValue, &config)
	defer jsonfile.Close()
	return config
}

func WatchFile(filePath string, change chan bool) {
	previousModTime := time.Now()

	for {
		time.Sleep(time.Second)

		file, _ := os.Stat(filePath)

		modTime := file.ModTime()
		if modTime.After(previousModTime) {
			change <- true
			previousModTime = modTime
		} else {
			change <- false
		}
	}
}
