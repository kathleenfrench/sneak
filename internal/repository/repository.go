package repository

import (
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/timshannon/bolthold"
)

// BoxRepository manages methods for working with box data
type BoxRepository interface {
	Save(box entity.Box) error
	Delete(id uint64) error
	GetByName(name string) (*entity.Box, error)
	Query(query *bolthold.Query) ([]entity.Box, error)
	List() []string
	BatchSave(boxes []entity.Box) error
	GetAll() ([]entity.Box, error)
}
