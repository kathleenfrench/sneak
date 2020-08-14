package pipeline

import (
	"github.com/kathleenfrench/common/fs"
)

// Manager interface wraps methods used for managing pipelines in sneak
type Manager interface {
	Save(p *Pipeline, path string) error
	CreateIfNotExist(path string) error
}

type pipelineManager struct {
	fs.Manager
}

// NewPipelineManager instantiates a new instance of the Manager interface
func NewPipelineManager() Manager {
	return &pipelineManager{}
}
