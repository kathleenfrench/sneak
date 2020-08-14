package pipeline

import (
	"fmt"

	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/repository"
	"github.com/kathleenfrench/sneak/pkg/file"
	"gopkg.in/yaml.v2"
)

type pipelineRepository struct {
	file         file.Manager
	manifestPath string
}

// NewPipelineRepository instantiates a new pipeline repository interface
func NewPipelineRepository(manifestPath string) repository.PipelineRepository {
	return &pipelineRepository{
		file:         file.NewManager(),
		manifestPath: manifestPath,
	}
}

func (r *pipelineRepository) ParseManifest() (*entity.PipelinesManifest, error) {
	return r.read()
}

func (r *pipelineRepository) SavePipeline(p *entity.Pipeline) error {
	manifest, err := r.read()
	if err != nil {
		return err
	}

	if manifest.Pipelines[p.Name] != nil {
		manifest.Pipelines[p.Name] = p
	} else {
		manifest.Pipelines[p.Name] = &entity.Pipeline{}
		manifest.Pipelines[p.Name] = p
	}

	return r.write(manifest)
}

func (r *pipelineRepository) SaveManifest(m *entity.PipelinesManifest) error {
	return r.write(m)
}

func (r *pipelineRepository) RemovePipeline(name string) error {
	manifest, err := r.read()
	if err != nil {
		return err
	}

	delete(manifest.Pipelines, name)
	return nil
}

func (r *pipelineRepository) ManifestExists() (bool, error) {
	manifestExists, err := r.file.FileExists(r.manifestPath)
	if err != nil {
		return false, err
	}

	if manifestExists {
		return true, nil
	}

	return false, nil
}

// ----------------------- hellpers

func (r *pipelineRepository) read() (*entity.PipelinesManifest, error) {
	// create file if it does not exist
	err := r.file.Touch(r.manifestPath)
	if err != nil {
		return nil, fmt.Errorf("could not create a pipeline file at %s - %w", r.manifestPath, err)
	}

	manifest := &entity.PipelinesManifest{}
	manifestYAML, err := r.file.ReadFile(r.manifestPath)
	if err != nil {
		return manifest, fmt.Errorf("could not read the pipeline manifest file at %s - %w", r.manifestPath, err)
	}

	err = yaml.Unmarshal(manifestYAML, &manifest)
	if err != nil {
		return manifest, fmt.Errorf("could not unmarshal pipeline manifest - %w", err)
	}

	return manifest, nil
}

func (r *pipelineRepository) write(m *entity.PipelinesManifest) error {
	pipelineFile, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not marshal pipeline manifest: %w", err)
	}

	err = r.file.Write(r.manifestPath, pipelineFile)
	if err != nil {
		return fmt.Errorf("could not save pipeline manifest: %w", err)
	}

	return nil
}
