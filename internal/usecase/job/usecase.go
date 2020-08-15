package job

import (
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
	"github.com/kathleenfrench/sneak/pkg/file"
)

// Usecase is an interface for methods controlling jobs in pipelines
type Usecase interface {
	SaveJob(job *entity.Job) error
	RemoveJob(jobName string) error
}

type jobUsecase struct {
	file file.Manager
	pipeline.Usecase
	pipelineName string
}

// NewJobUsecase instantiates a new job usecase interface
func NewJobUsecase(u pipeline.Usecase) Usecase {
	return &jobUsecase{
		Usecase: u,
		file:    file.NewManager(),
	}
}

func (u *jobUsecase) SaveJob(job *entity.Job) error {
	return nil
}

func (u *jobUsecase) RemoveJob(jobName string) error {
	return nil
}
