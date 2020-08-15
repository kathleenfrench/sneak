package htb

import (
	"github.com/kathleenfrench/sneak/internal/usecase/action"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
)

// ActionsGUI manages methods and properties of the terminal GUI re: actions
type ActionsGUI struct {
	usecase action.Usecase
}

// NewActionsGUI instantiates a new ActionsGUI struct
func NewActionsGUI(u pipeline.Usecase) *ActionsGUI {
	return &ActionsGUI{
		usecase: action.NewActionUsecase(u),
	}
}
