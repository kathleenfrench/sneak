package pipeline

import (
	"fmt"
	"io/ioutil"

	"github.com/kathleenfrench/common/fs"
	"gopkg.in/yaml.v2"
)

// Save saves the workflow file with pipelines
func Save(w *Workflow, wpath string) error {
	wfile, err := yaml.Marshal(w)
	if err != nil {
		return fmt.Errorf("could not marshal workflow file: %w", err)
	}

	// create file if it does not exist
	err = fs.CreateFile(wpath)
	if err != nil {
		return fmt.Errorf("could not create a workflow file at %s - %w", wpath, err)
	}

	err = ioutil.WriteFile(wpath, wfile, 0600)
	if err != nil {
		return fmt.Errorf("could not save workflow file: %w", err)
	}

	return nil
}
