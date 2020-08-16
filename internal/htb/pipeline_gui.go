package htb

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
)

// PipelineGUI is an interface for methods handling pipelines
type PipelineGUI interface {
	SelectPipelineFromDropdown(pipelines entity.Pipelines) *entity.Pipeline
	SelectPipelineActionsDropdown(pipeline *entity.Pipeline, pipelines entity.Pipelines) error
	PromptUserForPipelineData() (*entity.Pipeline, error)
	DefaultPipelineDropdownHandler() error
	ListPipelinesDropdown() error
	AddPipelineDropdown() error
	ManageWordlistsDropdown() error
}

// PipelineGUI manages methods and properties of the terminal GUI re: pipelines
type pipelineGUI struct {
	usecase pipeline.Usecase
	jobsGUI *JobsGUI
}

// NewPipelineGUI instantiates a new pipeline gui struct
func NewPipelineGUI(u pipeline.Usecase) PipelineGUI {
	return &pipelineGUI{
		usecase: u,
		jobsGUI: NewJobsGUI(u),
	}
}

var (
	pipelineListAll = "list all pipelines"
	addNewPipeline  = "add new pipeline"
	manageWordlists = "manage wordlists"
)

var defaultPipelineOpts = []string{
	pipelineListAll,
	addNewPipeline,
	manageWordlists,
	quit,
}

// DefaultPipelineDropdownHandler handles the GUI for sneak pipelines without argument
func (pg *pipelineGUI) DefaultPipelineDropdownHandler() error {
	choice := gui.SelectPromptWithResponse("select one", defaultPipelineOpts, nil, true)
	switch choice {
	case pipelineListAll:
		return pg.ListPipelinesDropdown()
	case addNewPipeline:
		return pg.AddPipelineDropdown()
	case manageWordlists:
		color.Red("TODO")
	case quit:
		os.Exit(0)
	}
	return nil
}

// AddPipelineDropdown is a GUI dropdown for adding a new pipeline
func (pg *pipelineGUI) AddPipelineDropdown() error {
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

// ListPipelinesDropdown is a dropdown GUI for all pipelines
func (pg *pipelineGUI) ListPipelinesDropdown() error {
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

// ManageWordlistsDropdown is the dropdown for managing wordlists
func (pg *pipelineGUI) ManageWordlistsDropdown() error {
	return nil
}
