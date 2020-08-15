package action

import (
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
)

// Usecase is an interface for methods controlling actions in pipelines
type Usecase interface {
	SaveAction(action *entity.Action) error
	GetAll() ([]*entity.Action, error)
	GetByName(name string) (*entity.Action, error)
	RemoveAction(name string) error
}

type actionUsecase struct {
	pipeline.Usecase
}

// NewActionUsecase instantiates a new actions usecase interface
func NewActionUsecase(u pipeline.Usecase) Usecase {
	return &actionUsecase{u}
}

func (u *actionUsecase) SaveAction(action *entity.Action) error {
	return nil
}

func (u *actionUsecase) GetAll() ([]*entity.Action, error) {
	return nil, nil
}

func (u *actionUsecase) GetByName(name string) (*entity.Action, error) {
	return nil, nil
}

func (u *actionUsecase) RemoveAction(name string) error {
	return nil
}
