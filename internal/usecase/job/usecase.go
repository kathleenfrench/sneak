package job

import (
	"github.com/fatih/color"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
	"github.com/kathleenfrench/sneak/pkg/file"
)

// Usecase is an interface for methods controlling jobs in pipelines
type Usecase interface {
	SaveJob(job *entity.Job, pipelineName string) error
	RemoveJob(jobName string) error
	GetPipelineJobs(pipelineName string) (map[string]*entity.Job, error)
}

type jobUsecase struct {
	file file.Manager
	pipeline.Usecase
}

// NewJobUsecase instantiates a new job usecase interface
func NewJobUsecase(u pipeline.Usecase) Usecase {
	return &jobUsecase{
		Usecase: u,
		file:    file.NewManager(),
	}
}

func (u *jobUsecase) SaveJob(job *entity.Job, pipelineName string) error {
	color.Green("pipeline name: %s", pipelineName)
	pipeline, err := u.GetByName(pipelineName)
	if err != nil {
		return err
	}
	color.Green("pipeline: %v", pipeline)

	if pipeline.Jobs == nil {
		pipeline.Jobs = make(map[string]*entity.Job)
	}

	pipeline.Jobs[job.Name] = &entity.Job{}
	pipeline.Jobs[job.Name] = job

	color.Green("pipeline jobs: %v", pipeline.Jobs)

	err = u.SavePipeline(pipeline)
	if err != nil {
		return err
	}

	return nil
}

func (u *jobUsecase) RemoveJob(jobName string) error {
	return nil
}

func (u *jobUsecase) GetPipelineJobs(pipelineName string) (map[string]*entity.Job, error) {
	p, err := u.GetByName(pipelineName)
	if err != nil {
		return nil, err
	}

	return p.Jobs, nil
}
