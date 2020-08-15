package pipeline

import (
	"fmt"

	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/repository"
	"github.com/kathleenfrench/sneak/pkg/file"
)

// Usecase is an interface for methods controlling pipelines
type Usecase interface {
	NewPipeline(p *entity.Pipeline) error
	GetAll() (entity.Pipelines, error)
	GetByName(name string) (*entity.Pipeline, error)
	RemovePipeline(name string) error
	ManifestExists() (bool, error)
	NewManifest() error
	SavePipeline(p *entity.Pipeline) error
}

type pipelineUsecase struct {
	Repository repository.PipelineRepository
	file       file.Manager
}

// NewPipelineUsecase instantiates a new pipeline usecase interface
func NewPipelineUsecase(r repository.PipelineRepository) Usecase {
	return &pipelineUsecase{
		Repository: r,
		file:       file.NewManager(),
	}
}

func (u *pipelineUsecase) ManifestExists() (bool, error) {
	return u.Repository.ManifestExists()
}

func (u *pipelineUsecase) NewManifest() error {
	manifestDefaults := &entity.PipelinesManifest{
		Version:   "v1",
		Pipelines: make(entity.Pipelines),
	}

	err := u.Repository.SaveManifest(manifestDefaults)
	if err != nil {
		return err
	}

	return nil
}

func (u *pipelineUsecase) SavePipeline(p *entity.Pipeline) error {
	return u.Repository.SavePipeline(p)
}

func (u *pipelineUsecase) NewPipeline(p *entity.Pipeline) error {
	return u.Repository.SavePipeline(p)
}

func (u *pipelineUsecase) GetAll() (entity.Pipelines, error) {
	manifest, err := u.Repository.ParseManifest()
	if err != nil {
		return nil, err
	}

	return manifest.Pipelines, nil
}

func (u *pipelineUsecase) GetByName(name string) (*entity.Pipeline, error) {
	manifest, err := u.Repository.ParseManifest()
	if err != nil {
		return nil, err
	}

	if p, found := manifest.Pipelines[name]; found {
		return p, nil
	}

	return nil, fmt.Errorf("no pipeline found by the name %s", name)
}

func (u *pipelineUsecase) RemovePipeline(name string) error {
	return u.Repository.RemovePipeline(name)
}
