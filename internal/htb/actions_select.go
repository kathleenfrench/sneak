package htb

import (
	"fmt"
	"os"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
)

const (
	returnToOtherActions = "return to other actions"
	removeAction         = "remove this action"
	viewRunner           = "view runner"
)

var singleActionOpts = []string{
	editDescription,
	viewRunner,
	removeAction,
	returnToOtherActions,
	quit,
}

// SelectActionFromDropdown lists a collection of defined actions and offers a host of options for how to interact with them
func (ag *ActionsGUI) SelectActionFromDropdown(actions map[string]*entity.Action) *entity.Action {
	names := getActionNames(actions)
	selection := gui.SelectPromptWithResponse("select an action", names, nil, false)
	selected := actions[selection]
	return selected
}

// SelectIndividualActionsActionsDropdown lists available actions for interacting/configuring/modifying sneak actions defined in the pipeline manifest
func (ag *ActionsGUI) SelectIndividualActionsActionsDropdown(a *entity.Action) error {
	actionChoice := gui.SelectPromptWithResponse("select from dropdown", singleActionOpts, nil, true)

	switch actionChoice {
	case editDescription:
		a.Description = gui.InputPromptWithResponse("provide a new description", a.Description, true)
		err := ag.usecase.SaveAction(a)
		if err != nil {
			return err
		}

		return ag.SelectIndividualActionsActionsDropdown(a)
	case viewRunner:
	case removeAction:
		var err error
		confirmRemoval := gui.ConfirmPrompt(fmt.Sprintf("are you sure you want to remove the action %s?", a.Name), "", false, true)
		switch confirmRemoval {
		case true:
			err = ag.usecase.RemoveAction(a.Name)
			if err != nil {
				return err
			}
		default:
			break
		}

		updatedActions, err := ag.usecase.GetAll()
		if err != nil {
			return err
		}

		newActionChoice := ag.SelectActionFromDropdown(updatedActions)
		return ag.SelectIndividualActionsActionsDropdown(newActionChoice)
	case returnToOtherActions:
		all, err := ag.usecase.GetAll()
		if err != nil {
			return err
		}

		return ag.SelectIndividualActionsActionsDropdown(ag.SelectActionFromDropdown(all))
	case quit:
		os.Exit(0)
	}

	return nil
}
