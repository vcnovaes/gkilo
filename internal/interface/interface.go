package tinterface

import (
	"common"
	"fmt"
)

func EditorRefreshScreen() {
    fmt.Print("\x1b[2J") // Clears the screen
    fmt.Print("\x1b[H")  // Moves the cursor to the top-left corner
    
    DrawRows()
    fmt.Print("\x1b[H");
}


func DrawRows() {
    _, h,_ := common.GetWindowSize()
    for y:=0; y< h;y++ {
        fmt.Print("~\r\n")
    }
}