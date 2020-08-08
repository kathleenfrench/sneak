package store

import (
	"fmt"

	"github.com/kathleenfrench/common/fs"
	"github.com/mitchellh/go-homedir"
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
