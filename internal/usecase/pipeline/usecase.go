package pipeline

import (
	"fmt"

	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/repository"
	"github.com/kathleenfrench/sneak/pkg/file"
	"gopkg.in/yaml.v2"
)

// Usecase is an interface for methods controlling pipelines
type Usecase interface {
	Save(p *entity.Pipeline) error
	GetAll() (*entity.Pipelines, error)
	GetByName(name string) (*entity.Pipeline, error)
	Remove(name string) error
}

type pipelineUsecase struct {
	Repository repository.PipelineRepository
	path       string
	file       file.Manager
}

// NewPipelineUsecase instantiates a new pipeline usecase interface
func NewPipelineUsecase(r repository.PipelineRepository, path string) Usecase {
	return &pipelineUsecase{
		Repository: r,
		path:       path,
		file:       file.NewManager(),
	}
}

func (u *pipelineUsecase) Save(p *entity.Pipeline) error {
	path := u.path
	pipelineFile, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("could not marshal pipeline file: %w", err)
	}

	// create file if it does not exist
	err = u.file.Touch(path)
	if err != nil {
		return fmt.Errorf("could not create a pipeline file at %s - %w", path, err)
	}

	err = u.file.Write(path, pipelineFile)
	if err != nil {
		return fmt.Errorf("could not save pipeline file: %w", err)
	}

	return nil
}

func (u *pipelineUsecase) GetAll() (*entity.Pipelines, error) {
	return nil, nil
}

func (u *pipelineUsecase) GetByName(name string) (*entity.Pipeline, error) {
	return nil, nil
}

func (u *pipelineUsecase) Remove(name string) error {
	return nil
}
