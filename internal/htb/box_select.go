package htb

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/helpers"
	"github.com/spf13/viper"
)

const (
	toggleActiveStatus = "toggle active status"
	checkConnection    = "check connection"
	openNotes          = "open notes editor"
	quickViewNotes     = "quickview notes"
	editDescription    = "edit description"
	flags              = "flags"
	returnToBoxes      = "return to other boxes"
	quit               = "quit"
	seeTable           = "show info table"
)

var boxActions = []string{
	toggleActiveStatus,
	checkConnection,
	openNotes,
	quickViewNotes,
	editDescription,
	flags,
	returnToBoxes,
	quit,
}

// SelectBoxActionsDropdown lists available actions with a single box or the ability to return to the 'main menu' of boxes
func (bg *BoxGUI) SelectBoxActionsDropdown(box entity.Box, boxes []entity.Box) error {
	if box.Active {
		bg.activeBox = box.Name
	}

	if !bg.singleBoxTableShown {
		PrintBoxDataTable(box)
		boxActions = append([]string{seeTable}, boxActions...)
	}

	bg.singleBoxTableShown = true
	selection := gui.SelectPromptWithResponse("select from the dropdown", boxActions, nil, true)

	switch selection {
	case seeTable:
		PrintBoxDataTable(box)
		return bg.SelectBoxActionsDropdown(box, boxes)
	case toggleActiveStatus:
		switch box.Active {
		case true:
			color.HiGreen("%s is currently set to active", box.Name)
			setInactive := gui.ConfirmPrompt(fmt.Sprintf("set %s as inactive?", box.Name), "", true, true)
			switch setInactive {
			case true:
				box.Active = false
				bg.activeBox = ""
			default:
				return bg.SelectBoxActionsDropdown(box, boxes)
			}
		default:
			color.Red("%s is not currently set to active", box.Name)
			setActive := gui.ConfirmPrompt(fmt.Sprintf("set %s as active?", box.Name), "", true, true)
			switch setActive {
			case true:
				box.Active = true
				bg.activeBox = box.Name
			default:
				return bg.SelectBoxActionsDropdown(box, boxes)
			}
		}

		// write the change to the db
		err := bg.usecase.Save(box)
		if err != nil {
			return err
		}

		color.Green("successfully changed the active status of %s!", box.Name)

		// after making that change, re-fetch all of our boxes for up to date info
		boxes, err = bg.usecase.GetAll()
		if err != nil {
			return err
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
	case quickViewNotes:
		note, err := checkForNoteFile(box.Name)
		if err != nil {
			return err
		}

		if len(note) == 0 {
			color.Yellow("you have not started a note for %s yet!", box.Name)
		} else {
			fmt.Println(helpers.RenderMarkdown(note))
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
	case checkConnection:
		err := helpers.SudoPing(box.IP)
		if err != nil {
			gui.Warn("uh oh, that box couldn't be reached! verify that the machine is active and your VPN connection is still intact", box.IP)
		} else {
			gui.Info("+1", "reachable!", box.IP)
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
	case openNotes:
		note, err := checkForNoteFile(box.Name)
		if err != nil {
			return err
		}

		updatedNote := gui.TextEditorInputAndSave(fmt.Sprintf("update your notes on %s in markdown", box.Name), note, viper.GetString("default_editor"))

		err = saveNoteFile(box.Name, updatedNote)
		if err != nil {
			return err
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
	case editDescription:
		box.Description = gui.InputPromptWithResponse("provide a new description", "", true)
		// write the change to the db
		err := bg.usecase.Save(box)
		if err != nil {
			return err
		}

		// after making that change, re-fetch all of our boxes for up to date info
		boxes, err = bg.usecase.GetAll()
		if err != nil {
			return err
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
	case flags:
		printFlagTable(box.Flags)
		addOrUpdate := gui.ConfirmPrompt("do you want to update any flag values?", "", false, true)
		switch addOrUpdate {
		case true:
			whichFlags := gui.SelectPromptWithResponse("which flag?", []string{"root", "user"}, nil, true)
			switch whichFlags {
			case "root":
				box.Flags.Root = gui.InputPromptWithResponse("enter a new root flag", box.Flags.Root, true)
			case "user":
				box.Flags.User = gui.InputPromptWithResponse("enter a new user flag", box.Flags.User, true)
			}

			// write the change to the db
			err := bg.usecase.Save(box)
			if err != nil {
				return err
			}

			// after making that change, re-fetch all of our boxes for up to date info
			boxes, err = bg.usecase.GetAll()
			if err != nil {
				return err
			}

			fallthrough
		default:
			return bg.SelectBoxActionsDropdown(box, boxes)
		}
	case returnToBoxes:
		return bg.SelectBoxActionsDropdown(bg.SelectBoxFromDropdown(boxes), boxes)
	case quit:
		os.Exit(0)
	}

	return nil
}

// SelectBoxFromDropdown lists a collection of boxes to choose from in a terminal dropdown
func (bg *BoxGUI) SelectBoxFromDropdown(boxes []entity.Box) entity.Box {
	boxNames, boxMap := makeGuiBoxMappings(boxes)
	selection := gui.SelectPromptWithResponse("select a box", boxNames, nil, false)
	selected := boxMap[selection]
	return selected
}
