package pipeline

import "github.com/kathleenfrench/sneak/pkg/file"

// Manager interface wraps methods used for managing pipelines in sneak
type Manager interface {
	Save(p *Pipelines) error
	GetAll() (*Pipelines, error)
	GetByName(name string) (*Pipeline, error)
}

type pipelineManager struct {
	file.Manager
	path string
}

// NewPipelineManager instantiates a new instance of the Manager interface
func NewPipelineManager(path string) Manager {
	return &pipelineManager{
		path: path,
	}
}
