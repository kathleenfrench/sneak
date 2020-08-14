package box

import (
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/repository"
	"github.com/timshannon/bolthold"
)

// Usecase contains methods for managing boxes
type Usecase interface {
	Save(box entity.Box) error
	Delete(id uint64) error
	GetByName(name string) (*entity.Box, error)
	Query(query *bolthold.Query) ([]entity.Box, error)
	List() []string
	BatchSave(boxes []entity.Box) error
	GetAll() ([]entity.Box, error)
}

type boxUsecase struct {
	Repository repository.BoxRepository
}

// NewUsecase instantiates a new box usecase interface
func NewUsecase(r repository.BoxRepository) Usecase {
	return &boxUsecase{
		Repository: r,
	}
}

func (u *boxUsecase) Save(box entity.Box) error {
	return u.Repository.Save(box)
}

func (u *boxUsecase) Delete(id uint64) error {
	return u.Repository.Delete(id)
}

func (u *boxUsecase) GetByName(name string) (*entity.Box, error) {
	return u.Repository.GetByName(name)
}

func (u *boxUsecase) Query(query *bolthold.Query) ([]entity.Box, error) {
	return u.Repository.Query(query)
}

func (u *boxUsecase) List() []string {
	return u.Repository.List()
}

func (u *boxUsecase) BatchSave(boxes []entity.Box) error {
	return u.Repository.BatchSave(boxes)
}

func (u *boxUsecase) GetAll() ([]entity.Box, error) {
	return u.Repository.GetAll()
}
