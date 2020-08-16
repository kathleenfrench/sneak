package htb

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/action"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
)

// PipelineGUI is an interface for methods handling pipelines
type PipelineGUI interface {
	SelectPipelineFromDropdown(pipelines entity.Pipelines) *entity.Pipeline
	SelectPipelineActionsDropdown(pipeline *entity.Pipeline, pipelines entity.Pipelines) error
	PromptUserForPipelineData() (*entity.Pipeline, error)
	ListPipelines() error
	AddPipeline() error
	DefaultPipelineDropdownHandler() error
}

// PipelineGUI manages methods and properties of the terminal GUI re: pipelines
type pipelineGUI struct {
	usecase       pipeline.Usecase
	actionUsecase action.Usecase
	jobsGUI       *JobsGUI
}

// NewPipelineGUI instantiates a new pipeline gui struct
func NewPipelineGUI(u pipeline.Usecase) PipelineGUI {
	return &pipelineGUI{
		usecase: u,
		jobsGUI: NewJobsGUI(u),
	}
}

var (
	wordlistsSub        = "wordlists"
	listPipelines       = "list all"
	viewPipelineActions = "actions"
	addNewPipeline      = "add new pipeline"
)

var defaultPipelineOpts = []string{
	listPipelines,
	addNewPipeline,
	viewPipelineActions,
	wordlistsSub,
	quit,
}

// DefaultPipelineDropdownHandler handles the GUI dropdown when no arguments are invoked on sneak pipeline
func (pg *pipelineGUI) DefaultPipelineDropdownHandler() error {
	choice := gui.SelectPromptWithResponse("select one", defaultPipelineOpts, nil, true)
	switch choice {
	case listPipelines:
		return pg.ListPipelines()
	case addNewPipeline:
		return pg.AddPipeline()
	case viewPipelineActions:
		return errors.New("actions")
	case wordlistsSub:
		color.Red("TODO")
	case quit:
		os.Exit(0)
	}
	return nil
}

// AddPipeline prompts the pipeline GUI for adding a new pipeline
func (pg *pipelineGUI) AddPipeline() error {
	newPipeline, err := pg.PromptUserForPipelineData()
	if err != nil {
		return err
	}

	err = pg.usecase.SavePipeline(newPipeline)
	if err != nil {
		return err
	}

	gui.Info("+1", fmt.Sprintf("%s was added successfully!", newPipeline.Name), newPipeline.Name)
	return nil
}

// ListPipelines prompts the pipeline GUI for choosing from all pipelines
func (pg *pipelineGUI) ListPipelines() error {
	pipelines, err := pg.usecase.GetAll()
	if err != nil {
		return err
	}

	if len(pipelines) == 0 {
		gui.Warn("you don't have any pipelines configured yet! run `sneak pipeline new` to get started", nil)
		return nil
	}

	selection := pg.SelectPipelineFromDropdown(pipelines)
	if err = pg.SelectPipelineActionsDropdown(selection, pipelines); err != nil {
		return err
	}

	return nil
}
