package htb

import (
	"os"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/job"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
)

// JobsGUI manages the GUI for interacting with jobs within a single pipeline
type JobsGUI struct {
	usecase   job.Usecase
	pipeline  *entity.Pipeline
	pipelines entity.Pipelines
	PipelineGUI
}

// NewJobsGUI instantiates a new JobsGUI
func NewJobsGUI(u pipeline.Usecase) *JobsGUI {
	return &JobsGUI{
		usecase: job.NewJobUsecase(u),
	}
}

// HandleJobsDropdown handles the initial jobs GUI dropdown when selectd within a single pipeline
func (jg *JobsGUI) HandleJobsDropdown(jobs map[string]*entity.Job) error {
	jobAction := gui.SelectPromptWithResponse("select from dropdown", []string{viewAllJobs, addNewJob, returnToPipeline, quit}, nil, true)

	switch jobAction {
	case viewAllJobs:
	case addNewJob:
	case returnToPipeline:
		return jg.SelectPipelineActionsDropdown(jg.pipeline, jg.pipelines)
	case quit:
		os.Exit(0)
	}

	return nil
}
