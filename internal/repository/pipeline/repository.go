package pipeline

import (
	"github.com/kathleenfrench/sneak/internal/repository"
	"github.com/timshannon/bolthold"
)

type pipelineRepository struct {
	*bolthold.Store
}

// NewPipelineRepository instantiates a new pipeline repository interface
func NewPipelineRepository(db *bolthold.Store) repository.PipelineRepository {
	return &pipelineRepository{db}
}
