package pipeline

import (
	"fmt"
	"io/ioutil"

	"github.com/kathleenfrench/common/fs"
	"gopkg.in/yaml.v2"
)

// Save saves the pipeline file with pipelines
func (m *pipelineManager) Save(p *Pipeline, path string) error {
	wfile, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("could not marshal pipeline file: %w", err)
	}

	// create file if it does not exist
	err = fs.CreateFile(path)
	if err != nil {
		return fmt.Errorf("could not create a pipeline file at %s - %w", path, err)
	}

	err = ioutil.WriteFile(path, wfile, 0600)
	if err != nil {
		return fmt.Errorf("could not save pipeline file: %w", err)
	}

	return nil
}

func (m *pipelineManager) CreateIfNotExist(path string) error {
	return nil
}
