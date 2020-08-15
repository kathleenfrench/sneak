package htb

import (
	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
)

const (
	viewAllJobs      = "see all jobs in pipeline"
	addNewJob        = "add new job"
	returnToPipeline = "return to pipeline"
)

const (
	returnToOtherJobs = "return to other jobs"
	disableJob        = "disable this job in the pipeline"
	removeJob         = "remove this job from the pipeline"
)

// SelectJobActionDropdown lists available actions for interacting/configuring an individual job
func (pg *JobsGUI) SelectJobActionDropdown(job *entity.Job) error {
	color.Green("SELECT JOB ACTIONS")
	return nil
}

// SelectJobFromDropdown lists a collection of jobs defined within a single pipeline where a job represents a single operation/script/task
func (pg *JobsGUI) SelectJobFromDropdown(jobs map[string]*entity.Job) *entity.Job {
	jobsNames := getJobKeys(jobs)
	selection := gui.SelectPromptWithResponse("select a job", jobsNames, nil, false)
	selected := jobs[selection]
	return selected
}
