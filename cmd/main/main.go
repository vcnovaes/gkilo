package main

import (
	"core"
)

func main() {
	var editor core.Editor

	editor.Init()
	editor.Run()
	defer editor.Close()

}
