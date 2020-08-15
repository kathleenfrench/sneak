package htb

import (
	"github.com/kathleenfrench/sneak/internal/usecase/action"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
)

// ActionsGUI is an interface for methods to work with actions within pipelines
type ActionsGUI interface {
}

// actionsGUI manages methods and properties of the terminal GUI re: actions
type actionsGUI struct {
	usecase action.Usecase
}

// NewActionsGUI instantiates a new ActionsGUI struct
func NewActionsGUI(u pipeline.Usecase) ActionsGUI {
	return &actionsGUI{
		usecase: action.NewActionUsecase(u),
	}
}
