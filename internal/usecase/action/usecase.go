package action

import (
	"fmt"

	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
)

// Usecase is an interface for methods controlling actions in pipelines
type Usecase interface {
	SaveAction(action *entity.Action) error
	GetAll() (map[string]*entity.Action, error)
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

func (u *actionUsecase) GetAll() (map[string]*entity.Action, error) {
	manifest, err := u.GetManifest()
	if err != nil {
		return nil, err
	}

	if manifest.Actions == nil {
		manifest.Actions = make(map[string]*entity.Action)
	}

	return manifest.Actions, nil
}

func (u *actionUsecase) GetByName(name string) (*entity.Action, error) {
	manifest, err := u.GetManifest()
	if err != nil {
		return nil, err
	}

	if manifest.Actions == nil {
		manifest.Actions = make(map[string]*entity.Action)
	}

	if m, ok := manifest.Actions[name]; ok {
		return m, nil
	}

	return nil, fmt.Errorf("action %q not found", name)
}

func (u *actionUsecase) RemoveAction(name string) error {
	_, err := u.GetByName(name)
	if err != nil {
		return err
	}

	manifest, err := u.GetManifest()
	if err != nil {
		return err
	}

	delete(manifest.Actions, name)
	return u.SaveManifest(manifest)
}
