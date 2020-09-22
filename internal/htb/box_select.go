package htb

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/action"
	"github.com/kathleenfrench/sneak/pkg/utils"
	"github.com/spf13/viper"
)

// shared
const (
	editDescription = "edit description"
	quit            = "quit"
)

const (
	toggleActiveStatus = "toggle active status"
	checkConnection    = "check connection"
	openNotes          = "open notes editor"
	quickViewNotes     = "quickview notes"
	notes              = "notes"
	flags              = "flags"
	returnToBoxes      = "return to other boxes"
	seeTable           = "show info table"
	runPipeline        = "run pipeline"
	runOneoffAction    = "run one-off action"
)

var boxActions = []string{
	toggleActiveStatus,
	checkConnection,
	runPipeline,
	runOneoffAction,
	notes,
	editDescription,
	flags,
	returnToBoxes,
	quit,
}

var notesActions = []string{
	openNotes,
	quickViewNotes,
}

// SelectBoxActionsDropdown lists available actions with a single box or the ability to return to the 'main menu' of boxes
func (bg *BoxGUI) SelectBoxActionsDropdown(box entity.Box, boxes []entity.Box) error {
	if box.Active {
		bg.activeBox = box
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
	case notes:
		notesNext := gui.SelectPromptWithResponse("select one", notesActions, nil, true)
		switch notesNext {
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
		case quickViewNotes:
			note, err := checkForNoteFile(box.Name)
			if err != nil {
				return err
			}

			if len(note) == 0 {
				color.Yellow("you have not started a note for %s yet!", box.Name)
			} else {
				fmt.Println(utils.RenderMarkdown(note))
			}

			return bg.SelectBoxActionsDropdown(box, boxes)
		}
	case runPipeline:
		empty := entity.Box{}
		if bg.activeBox == empty {
			gui.Warn("you must set a box as active before running a pipeline", box.Name)
			return bg.SelectBoxActionsDropdown(box, boxes)
		}

		all, err := bg.pipUsecase.GetAll()
		if err != nil {
			return err
		}

		lsNamesMap, names := genPipelineRunnerNameListMap(all)
		run := gui.SelectPromptWithResponse("select a pipeline to run", names, nil, false)
		chosen := lsNamesMap[run]
		toRun := all[chosen]
		return bg.RunPipeline(toRun)
	case runOneoffAction:
		actionUsecase := action.NewActionUsecase(bg.pipUsecase)
		allActions, err := actionUsecase.GetAll()
		if err != nil {
			return err
		}

		oneoff := gui.SelectPromptWithResponse("select a one-off action", getActionNames(allActions), nil, false)
		chosenOneoff := allActions[oneoff]
		return bg.HandleRunnerAction(chosenOneoff)
	case toggleActiveStatus:
		switch box.Active {
		case true:
			color.HiGreen("%s is currently set to active", box.Name)
			setInactive := gui.ConfirmPrompt(fmt.Sprintf("set %s as inactive?", box.Name), "", true, true)
			switch setInactive {
			case true:
				box.Active = false
				bg.activeBox = box
			default:
				return bg.SelectBoxActionsDropdown(box, boxes)
			}
		default:
			color.Red("%s is not currently set to active", box.Name)
			setActive := gui.ConfirmPrompt(fmt.Sprintf("set %s as active?", box.Name), "", true, true)
			switch setActive {
			case true:
				box.Active = true
				bg.activeBox = box
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
	case checkConnection:
		err := utils.SudoPing(box.IP)
		if err != nil {
			gui.Warn("uh oh, that box couldn't be reached! verify that the machine is active and your VPN connection is still intact", box.IP)
		} else {
			gui.Info("+1", "reachable!", box.IP)
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

func genPipelineRunnerNameListMap(pipelines map[string]*entity.Pipeline) (map[string]string, []string) {
	namesMap := make(map[string]string)
	names := []string{}

	for n, p := range pipelines {
		name := fmt.Sprintf("[%s]: %s", color.YellowString(n), p.Description)
		namesMap[name] = n
		names = append(names, name)
	}

	return namesMap, names
}
