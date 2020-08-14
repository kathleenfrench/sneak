package htb

import (
	boxusecase "github.com/kathleenfrench/sneak/internal/usecase/box"
)

// BoxGUI is a struct for managing the box GUI
type BoxGUI struct {
	singleBoxTableShown bool
	activeBox           string
	usecase             boxusecase.Usecase
}

// NewBoxGUI instantiates a new box gui interface
func NewBoxGUI(use boxusecase.Usecase) *BoxGUI {
	return &BoxGUI{
		singleBoxTableShown: false,
		activeBox:           "",
		usecase:             use,
	}
}
