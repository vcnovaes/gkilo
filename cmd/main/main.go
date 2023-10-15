package main

import (
	"common"
	tinterface "interface"
	"os"
	"syscall"
	"unsafe"
)

type EditorConfig struct {
    originTermios *syscall.Termios
    screenRows int 
    screenCols int
}



const CTRL_Q byte = 17

func getTermios(fd uintptr) *syscall.Termios {
    var t syscall.Termios
    _, _, err := syscall.Syscall6(
        syscall.SYS_IOCTL, // Input output control
        os.Stdin.Fd(),
        syscall.TCGETS,
        uintptr(unsafe.Pointer(&t)),
        0, 0, 0)

    if err != 0 {
        panic("Error getting termios")
    }

    return &t
}

func setTermios(fd uintptr, term *syscall.Termios) {
    _, _, err := syscall.Syscall6(
        syscall.SYS_IOCTL,
        os.Stdin.Fd(),
        syscall.TCSETS,
        uintptr(unsafe.Pointer(term)),
        0,0,0)
    if err != 0 {
        panic("err"); 
    }
}

func setRaw(term *syscall.Termios) {
    // This attempts to replicate the behaviour documented for cfmakeraw in
    // the termios(3) manpage.
    term.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
    // newState.Oflag &^= syscall.OPOST
    term.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
    term.Cflag &^= syscall.CSIZE | syscall.PARENB
    term.Cflag |= syscall.CS8

    term.Cc[syscall.VMIN] = 1
    term.Cc[syscall.VTIME] = 0
}


func setupRawMode(editorConfig *EditorConfig) {
    editorConfig.originTermios = getTermios(os.Stdin.Fd())
	defer setTermios(
			os.Stdin.Fd(), 
			editorConfig.originTermios)
	
	setRaw(editorConfig.originTermios)
	setTermios(os.Stdin.Fd(), editorConfig.originTermios)
}

func runEditor() {
    for {
        tinterface.EditorRefreshScreen()
        processKeyPressed()
    }
}

func processKeyPressed() {
	buffer := make([]byte, 1)
	syscall.Read(0, buffer)
	ch := string(buffer)
	print(ch)
    
	if buffer[0] == CTRL_Q {
		tinterface.EditorRefreshScreen()
        os.Exit(0)
	}
}

func initEditor(editorConfig *EditorConfig) {
    editorConfig.screenCols, editorConfig.screenRows, _ = common.GetWindowSize()
    setupRawMode(editorConfig)
}


func main() {
    var editorConfig EditorConfig
    initEditor(&editorConfig)
    runEditor()
}