package store

import (
	"fmt"
	"time"

	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

// Config represent config values for the bolt db
type Config struct {
	FileName string
	Path     string
}

func (c *Config) constructFullPath() string {
	return fmt.Sprintf("%s/%s", c.Path, c.FileName)
}

// NewDB initializes a new DB
func NewDB(c *Config) (*bolthold.Store, error) {
	opts := &bolthold.Options{
		Options: &bbolt.Options{
			Timeout: 5 * time.Second,
		},
	}

	db, err := bolthold.Open(c.constructFullPath(), 0600, opts)
	if err != nil {
		return db, err
	}

	return db, nil
}
