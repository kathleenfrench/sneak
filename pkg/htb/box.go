package htb

import (
	"fmt"
	"time"

	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

// Box represents a machine
type Box struct {
	Name        string `boldholdKey:"Name"`
	IP          string
	Completed   bool `boltholdIndex:"Completed"` // when root + user flags captured
	Active      bool `boltholdIndex:"Active"`
	Hostname    string
	OS          string
	Difficulty  string // easy, medium, hard, insane
	Notes       string
	Up          bool
	Flags       Flags
	Created     time.Time
	LastUpdated time.Time
}

// Flags represents the capture flags on a box
type Flags struct {
	Root string
	User string
}

// bucketName returns the bucket name for box data
func (bx *Box) bucketName() string {
	return "Boxes"
}

// CreateBox inserts a new box into the db
func CreateBox(db *bolthold.Store, box Box) error {
	if box.Created.IsZero() {
		box.Created = time.Now()
		box.LastUpdated = box.Created
	}

	err := db.Upsert(box.Name, box)
	if err != nil {
		return err
	}

	return nil
	// if err := db.Bolt().Update(func(tx *bbolt.Tx) error {
	// 	if box.Created.IsZero() {
	// 		box.Created = time.Now()
	// 		box.LastUpdated = box.Created
	// 	}

	// 	if box.ID != 0 {
	// 		var incrErr error
	// 		bucket := tx.Bucket([]byte(box.bucketName()))
	// 		box.ID, incrErr = bucket.NextSequence()
	// 		if incrErr != nil {
	// 			gui.Warn(fmt.Sprintf("issue autoincrementing: %v", incrErr), nil)
	// 		}

	// 		box.Created = time.Now()
	// 	} else {
	// 		box.LastUpdated = time.Now()
	// 	}

	// 	color.HiGreen("next int: %d", box.ID)

	// 	err := db.TxUpsert(tx, box.ID, box)
	// 	return err
	// }); err != nil {
	// 	return err
	// }

	// return nil
}

// DeleteBox deletes a box
func DeleteBox(db *bolthold.Store, boxID uint64) error {
	err := db.DeleteMatching(&Box{}, bolthold.Where(bolthold.Key).Eq(boxID))
	if err != nil {
		return err
	}

	return nil
}

// GetBoxByID gets a box by key
func GetBoxByID(db *bolthold.Store, id uint64) (*Box, error) {
	b := &Box{}
	err := db.Get(id, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetBoxByName fetches a box by name
func GetBoxByName(db *bolthold.Store, name string) (*Box, error) {
	b := &Box{}
	err := db.Find(&b, bolthold.Where(bolthold.Key).Eq(name))
	if err != nil {
		return nil, err
	}

	return b, nil
}

// QueryBoxes fetches a collection of boxes based off a query
func QueryBoxes(db *bolthold.Store, query *bolthold.Query) ([]Box, error) {
	var boxes []Box

	err := db.Find(&boxes, query)
	if err != nil {
		return nil, err
	}

	return boxes, nil
}

// GetAllBoxes returns all boxes in the db
func GetAllBoxes(db *bolthold.Store) ([]Box, error) {
	boxes := []Box{}
	db.Find(&boxes, bolthold.Where(bolthold.Key).Ne(""))
	return boxes, nil
}

// List returns a list of hte boxes
func (bx Box) List(db *bolthold.Store) []string {
	var all []string
	var boxes []Box

	if err := db.Bolt().View(func(tx *bbolt.Tx) error {
		return db.Find(&boxes, bolthold.Where("ID").Gt(uint64(0)))
	}); err != nil {
		panic(err)
	} else {
		for _, r := range boxes {
			all = append(all, fmt.Sprintf("%v", r))
		}
	}

	return all
}

// AddBoxes batch inserts boxes to the database
func AddBoxes(db *bolthold.Store, boxes []*Box) error {
	return db.Bolt().Update(func(tx *bbolt.Tx) error {
		for i := range boxes {
			err := db.TxInsert(tx, boxes[i].Name, boxes[i])
			if err != nil {
				return err
			}
		}

		all := &bolthold.Query{}
		return db.TxDeleteMatching(tx, &Box{}, all.SortBy("Name"))
	})
}
