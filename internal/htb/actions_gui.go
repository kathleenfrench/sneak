package htb

import (
	"os"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/action"
)

// ActionsGUI is an interface for methods to work with actions within pipelines
type ActionsGUI struct {
	usecase action.Usecase
	PipelineGUI
}

// NewActionsGUI instantiates a new ActionsGUI struct
func NewActionsGUI(u action.Usecase) *ActionsGUI {
	return &ActionsGUI{
		usecase: u,
	}
}

const (
	viewAllActions string = "view all actions"
	addNewAction   string = "add new action"
)

var actionDropdownOpts = []string{
	viewAllActions,
	addNewAction,
	quit,
}

// HandleActionsDropdown handles the initial dropdown for the actions GUI when selecting next steps for how to interact with one's defined (or undefined) actions
func (ag *ActionsGUI) HandleActionsDropdown(actions map[string]*entity.Action) error {
	actionSelect := gui.SelectPromptWithResponse("select from dropdown", actionDropdownOpts, nil, true)

	switch actionSelect {
	case viewAllActions:
		if actions == nil {
			gui.Warn("you do not have any actions defined yet", nil)
			return ag.HandleActionsDropdown(actions)
		}

		selected := ag.SelectActionFromDropdown(actions)
		return ag.SelectIndividualActionsActionsDropdown(selected)
	case addNewAction:
		newAction := &entity.Action{
			Name:        gui.InputPromptWithResponse("what do you want to call this action? (no spaces)", "", true),
			Description: gui.InputPromptWithResponse("describe what this action does", "", true),
		}

		// TODO: ADD RUNNER

		err := ag.usecase.SaveAction(newAction)
		if err != nil {
			return err
		}

		updatedActions, err := ag.usecase.GetAll()
		if err != nil {
			return err
		}

		return ag.HandleActionsDropdown(updatedActions)
	case quit:
		os.Exit(0)
	}

	return nil
}

func getActionNames(actions map[string]*entity.Action) []string {
	names := []string{}
	for n := range actions {
		names = append(names, n)
	}

	return names
}
