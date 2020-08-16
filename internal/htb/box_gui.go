package htb

import (
	"fmt"
	"os"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/usecase/box"
)

// BoxGUI is a struct for managing the box GUI
type BoxGUI struct {
	singleBoxTableShown bool
	activeBox           string
	usecase             box.Usecase
}

// NewBoxGUI instantiates a new box gui struct
func NewBoxGUI(use box.Usecase) *BoxGUI {
	return &BoxGUI{
		singleBoxTableShown: false,
		activeBox:           "",
		usecase:             use,
	}
}

var (
	listBoxes = "list"
	newBox    = "new"
)

var defaultBoxOpts = []string{
	listBoxes,
	newBox,
	quit,
}

// DefaultDropdownHandler is the dropdown GUI for when `sneak box` is run without arguments
func (bg *BoxGUI) DefaultDropdownHandler() error {
	choice := gui.SelectPromptWithResponse("select one", defaultBoxOpts, nil, true)
	switch choice {
	case listBoxes:
		boxes, err := bg.usecase.GetAll()
		if err != nil {
			return err
		}

		if len(boxes) == 0 {
			gui.Warn("you don't have any boxes yet! run `sneak box new` to get started", nil)
			return nil
		}

		selection := bg.SelectBoxFromDropdown(boxes)

		if err = bg.SelectBoxActionsDropdown(selection, boxes); err != nil {
			return err
		}
	case newBox:
		box, err := bg.PromptUserForBoxData()
		if err != nil {
			return err
		}

		err = bg.usecase.Save(box)
		if err != nil {
			return err
		}

		gui.Info("+1", fmt.Sprintf("%s was added successfully!", box.Name), fmt.Sprintf("%s", box.IP))
	case quit:
		os.Exit(0)
	}

	return nil
}
