package htb

import (
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
)

// PipelineGUI is an interface for methods handling pipelines
type PipelineGUI interface {
	SelectPipelineFromDropdown(pipelines entity.Pipelines) *entity.Pipeline
	SelectPipelineActionsDropdown(pipeline *entity.Pipeline, pipelines entity.Pipelines) error
	PromptUserForPipelineData() (*entity.Pipeline, error)
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
