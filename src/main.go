package main

import (
	"github.com/vcnovaes/gkilo/src/internal/core"
)

func main() {
	var editor core.Editor

	editor.Init()
	editor.Run()
	defer editor.Close()

}
