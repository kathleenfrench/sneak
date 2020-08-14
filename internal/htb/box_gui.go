package htb

import (
	"github.com/kathleenfrench/sneak/internal/usecase/box"
)

// BoxGUI is a struct for managing the box GUI
type BoxGUI struct {
	singleBoxTableShown bool
	activeBox           string
	usecase             box.Usecase
}

// NewBoxGUI instantiates a new box gui interface
func NewBoxGUI(use box.Usecase) *BoxGUI {
	return &BoxGUI{
		singleBoxTableShown: false,
		activeBox:           "",
		usecase:             use,
	}
}
