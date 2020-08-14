package box

import (
	"fmt"
	"time"

	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/repository"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

type boxRepository struct {
	*bolthold.Store
}

// NewBoxRepository instantiates a new box repository interface
func NewBoxRepository(db *bolthold.Store) repository.BoxRepository {
	return &boxRepository{db}
}

func (r *boxRepository) Save(box entity.Box) error {
	if box.Created.IsZero() {
		box.Created = time.Now()
		box.LastUpdated = box.Created
	} else {
		box.LastUpdated = time.Now()
	}

	err := r.Store.Upsert(box.Name, box)
	if err != nil {
		return err
	}

	return nil
}

func (r *boxRepository) Delete(id uint64) error {
	err := r.Store.DeleteMatching(&entity.Box{}, bolthold.Where(bolthold.Key).Eq(id))
	if err != nil {
		return err
	}

	return nil
}

func (r *boxRepository) GetByName(name string) (*entity.Box, error) {
	b := &entity.Box{}
	err := r.Store.Find(&b, bolthold.Where(bolthold.Key).Eq(name))
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *boxRepository) Query(query *bolthold.Query) ([]entity.Box, error) {
	var boxes []entity.Box

	err := r.Store.Find(&boxes, query)
	if err != nil {
		return nil, err
	}

	return boxes, nil
}

func (r *boxRepository) GetAll() ([]entity.Box, error) {
	boxes := []entity.Box{}
	r.Store.Find(&boxes, bolthold.Where(bolthold.Key).Ne(""))
	return boxes, nil
}

func (r *boxRepository) List() []string {
	var all []string
	var boxes []*entity.Box

	if err := r.Store.Bolt().View(func(tx *bbolt.Tx) error {
		return r.Store.Find(&boxes, bolthold.Where("ID").Gt(uint64(0)))
	}); err != nil {
		panic(err)
	} else {
		for _, r := range boxes {
			all = append(all, fmt.Sprintf("%v", r))
		}
	}

	return all
}

func (r *boxRepository) BatchSave(boxes []entity.Box) error {
	return r.Store.Bolt().Update(func(tx *bbolt.Tx) error {
		for i := range boxes {
			err := r.Store.TxInsert(tx, boxes[i].Name, boxes[i])
			if err != nil {
				return err
			}
		}

		all := &bolthold.Query{}
		return r.Store.TxDeleteMatching(tx, &entity.Box{}, all.SortBy("Name"))
	})
}
