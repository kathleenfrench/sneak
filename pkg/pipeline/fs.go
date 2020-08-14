package pipeline

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Save saves the pipeline file with pipelines
func (m *pipelineManager) Save(p *Pipeline) error {
	path := m.path
	pipelineFile, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("could not marshal pipeline file: %w", err)
	}

	// create file if it does not exist
	err = m.Touch(path)
	if err != nil {
		return fmt.Errorf("could not create a pipeline file at %s - %w", path, err)
	}

	err = m.Write(path, pipelineFile)
	if err != nil {
		return fmt.Errorf("could not save pipeline file: %w", err)
	}

	return nil
}
