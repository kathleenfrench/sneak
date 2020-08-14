package store

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/sneak/internal/repository/box"
	"github.com/mitchellh/go-homedir"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

// GetDataDirectory parses the path to .sneak's expected data directory, checks if the directory exists, attempts to create it if it does not, then returns the path
func GetDataDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return home, err
	}

	dir := fmt.Sprintf("%s/.sneak", home)

	err = fs.CreateDir(dir)
	if err != nil {
		return dir, err
	}

	return dir, nil
}

var sneakBuckets = map[string]bool{
	"Box": true,
}

// Buckets returns all of the db buckets for sneak
func Buckets(db *bolthold.Store, name string) []string {
	buckets := []string{}

	if len(name) > 0 {
		name = name[1:]
	}

	switch name {
	case "Box":
		boxRepository := box.NewBoxRepository(db)
		buckets = append(buckets, strings.Join(boxRepository.List(), "\n"))
	default:
		buckets = append(buckets, fmt.Sprintf("%s is not a valid bucket in sneak", name))
	}

	return buckets
}

// EmptyBuckets resets buckets by name
func EmptyBuckets(db *bolthold.Store, name string) []string {
	out := []string{}
	switch name {
	case "/all":
		for b := range sneakBuckets {
			b = strings.Title(strings.ToLower(b))
			out = append(out, reset(db, b))
		}
	default:
		name = name[1:]
		if sneakBuckets[name] {
			out = append(out, reset(db, name))
		} else {
			out = append(out, "choose a bucket")
		}
	}

	return out
}

func reset(db *bolthold.Store, name string) string {
	if sneakBuckets[name] {
		if err := db.Bolt().Update(func(tx *bbolt.Tx) error {
			tx.DeleteBucket([]byte(name))
			_, err := tx.CreateBucket([]byte(name))
			return err
		}); err != nil {
			return fmt.Sprintf("%s Bucket --> %s", name, err.Error())
		}

		return fmt.Sprintf("%s Bucket --> Reset", name)
	}

	return fmt.Sprintf("%s Bucket --> Does not Exist", name)
}

// Backup backs up the current database to the specified directory
func Backup(dir string) error {
	dbName := "sneak.db"
	dbBackupName := "sneak_backup.db"

	src, err := os.Open(fmt.Sprintf("%s/%s", dir, dbName))
	if err != nil {
		return err
	}

	defer src.Close()

	backupFile, err := os.Create(fmt.Sprintf("%s/%s", dir, dbBackupName))
	if err != nil {
		return err
	}

	defer backupFile.Close()

	_, err = io.Copy(backupFile, src)
	if err != nil {
		return err
	}

	return nil
}
