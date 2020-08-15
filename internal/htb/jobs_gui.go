package htb

import (
	"os"

	"github.com/fatih/color"
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

var jobsDropdownOpts = []string{
	viewAllJobs,
	addNewJob,
	returnToPipeline,
	quit,
}

// HandleJobsDropdown handles the initial jobs GUI dropdown when selectd within a single pipeline
func (jg *JobsGUI) HandleJobsDropdown(jobs map[string]*entity.Job) error {
	jobAction := gui.SelectPromptWithResponse("select from dropdown", jobsDropdownOpts, nil, true)

	switch jobAction {
	case viewAllJobs:
		if jobs == nil {
			gui.Warn("you do not have any jobs set on this pipeline", jg.pipeline.Name)
			return jg.SelectPipelineActionsDropdown(jg.pipeline, jg.pipelines)
		}

		selectedJob := jg.SelectJobFromDropdown(jobs)
		return jg.SelectJobActionDropdown(selectedJob)
	case addNewJob:
		newJob := &entity.Job{
			Name:        gui.InputPromptWithResponse("what do you want to call this job? (no spaces)", "", true),
			Description: gui.InputPromptWithResponse("describe what this job does", "", true),
		}

		err := jg.usecase.SaveJob(newJob, jg.pipeline.Name)
		if err != nil {
			return err
		}

		jobs, err := jg.usecase.GetPipelineJobs(jg.pipeline.Name)
		if err != nil {
			return err
		}

		color.Yellow("JOBS: %v", jobs)

		return jg.HandleJobsDropdown(jobs)
	case returnToPipeline:
		return jg.SelectPipelineActionsDropdown(jg.pipeline, jg.pipelines)
	case quit:
		os.Exit(0)
	}

	return nil
}
