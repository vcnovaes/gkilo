package common

import (
	"os"

	"golang.org/x/term"
)
func GetWindowSize() (int, int, error) {
	fd := int(os.Stdout.Fd())
	w, h, err := term.GetSize(fd)
	if err != nil {
		return 0, 0, err
	}
	return w, h, err
}