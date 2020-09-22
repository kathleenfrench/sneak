package htb

import (
	"fmt"
	"os"

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
	manageActions     = "manage actions"
)

var singleJobActionsDropdown = []string{
	editDescription,
	manageActions,
	disableJob,
	removeJob,
	returnToOtherJobs,
	quit,
}

var (
	addActionToJob      = "add another action to job"
	viewExistingActions = "view existing job actions"
)

var singleJobActionOpts = []string{
	addActionToJob,
	viewExistingActions,
	quit,
}

// SelectJobActionDropdown lists available actions for interacting/configuring an individual job
func (jg *JobsGUI) SelectJobActionDropdown(job *entity.Job) error {
	jobAction := gui.SelectPromptWithResponse("select from dropdown", singleJobActionsDropdown, nil, true)
	switch jobAction {
	case editDescription:
		job.Description = gui.InputPromptWithResponse("provide a new description", job.Description, true)
		err := jg.usecase.SaveJob(job, jg.pipeline.Name)
		if err != nil {
			return err
		}

		return jg.SelectJobActionDropdown(job)
	case manageActions:
		choice := gui.SelectPromptWithResponse("select one", singleJobActionOpts, nil, false)
		switch choice {
		case addActionToJob:
			all, err := jg.actionsUsecase.GetAll()
			if err != nil {
				return err
			}

			toAdd := gui.SelectPromptWithResponse("select an action to add", getActionNames(all), nil, true)
			job.Actions = append(job.Actions, toAdd)
			err = jg.usecase.SaveJob(job, jg.pipeline.Name)
			if err != nil {
				return err
			}
		case viewExistingActions:
			if job.Actions == nil {
				gui.Warn("you do not have any actions defined for this job yet", nil)
				return jg.SelectJobActionDropdown(job)
			}

			jobActions, err := jg.actionsUsecase.GetJobActions(job.Actions)
			if err != nil {
				return err
			}

			agui := NewActionsGUI(jg.actionsUsecase)
			selected := agui.SelectActionFromDropdown(jobActions)
			return agui.SelectIndividualActionsActionsDropdown(selected)
		case quit:
			os.Exit(0)
		}
	case disableJob:
		var enabledStatus bool
		switch job.Disabled {
		case true:
			enabledStatus = gui.ConfirmPrompt(fmt.Sprintf("%s is currently disabled, re-enable it?", job.Name), "", true, true)
		default:
			enabledStatus = gui.ConfirmPrompt(fmt.Sprintf("%s is currently enabled, disable it?", job.Name), "", true, true)
		}

		job.Disabled = enabledStatus
		err := jg.usecase.SaveJob(job, jg.pipeline.Name)
		if err != nil {
			return err
		}

		return jg.SelectJobActionDropdown(job)
	case removeJob:
		var err error
		confirmRemoval := gui.ConfirmPrompt(fmt.Sprintf("are you sure you want to remove %s from the %s pipeline?", job.Name, jg.pipeline.Name), "", false, true)
		switch confirmRemoval {
		case true:
			err = jg.usecase.RemoveJob(job.Name, jg.pipeline.Name)
			if err != nil {
				return err
			}
		default:
			break
		}

		jg.pipeline.Jobs, err = jg.usecase.GetPipelineJobs(jg.pipeline.Name)
		if err != nil {
			return err
		}

		newJob := jg.SelectJobFromDropdown(jg.pipeline.Jobs)
		return jg.SelectJobActionDropdown(newJob)
	case returnToOtherJobs:
		newJob := jg.SelectJobFromDropdown(jg.pipeline.Jobs)
		return jg.SelectJobActionDropdown(newJob)
	case quit:
		os.Exit(0)
	}

	return nil
}

// SelectJobFromDropdown lists a collection of jobs defined within a single pipeline where a job represents a single operation/script/task
func (jg *JobsGUI) SelectJobFromDropdown(jobs map[string]*entity.Job) *entity.Job {
	printJobsTable(jobs)
	jobsNames := getJobKeys(jobs)
	selection := gui.SelectPromptWithResponse("select a job", jobsNames, nil, true)
	selected := jobs[selection]
	return selected
}
