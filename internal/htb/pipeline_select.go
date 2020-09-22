package htb

import (
	"os"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
)

const (
	returnToPipelines = "return to other pipelines"
	wordlists         = "wordlists"
	actions           = "actions"
	jobs              = "jobs"
)

var pipelineActions = []string{
	jobs,
	wordlists,
	editDescription,
	returnToPipelines,
	quit,
}

// SelectPipelineActionsDropdown lists available actions for interacting/configuring an individual pipeline
func (pg *pipelineGUI) SelectPipelineActionsDropdown(pipeline *entity.Pipeline, pipelines entity.Pipelines) error {
	selection := gui.SelectPromptWithResponse("select from dropdown", pipelineActions, nil, true)

	if pipeline.Jobs == nil {
		pipeline.Jobs = make(map[string]*entity.Job)
	}

	switch selection {
	case jobs:
		if pipeline.Jobs == nil {
			pipeline.Jobs = make(map[string]*entity.Job)
		}

		pg.jobsGUI.pipeline = pipeline
		pg.jobsGUI.pipelines = pipelines
		return pg.jobsGUI.HandleJobsDropdown(pipeline.Jobs)
	case wordlists:
		color.Red("TODO")
	case editDescription:
		pipeline.Description = gui.InputPromptWithResponse("provide a new description", "", true)
		err := pg.usecase.SavePipeline(pipeline)
		if err != nil {
			return err
		}

		pipelines, err = pg.usecase.GetAll()
		if err != nil {
			return err
		}

		return pg.SelectPipelineActionsDropdown(pipeline, pipelines)
	case returnToPipelines:
		return pg.SelectPipelineActionsDropdown(pg.SelectPipelineFromDropdown(pipelines), pipelines)
	case quit:
		os.Exit(0)
	}

	return nil
}

// SelectPipelineFromDropdown lists a collection of pipelines to choose from in a terminal dropdown
func (pg *pipelineGUI) SelectPipelineFromDropdown(pipelines entity.Pipelines) *entity.Pipeline {
	pipelineNames := getPipelineMapKeys(pipelines)
	selection := gui.SelectPromptWithResponse("select a pipeline", pipelineNames, nil, false)
	selected := pipelines[selection]
	return selected
}
