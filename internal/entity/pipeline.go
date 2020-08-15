package entity

import (
	"fmt"
	"strings"
)

// PipelinesManifest represents the physical file representation of pipeline data
type PipelinesManifest struct {
	Version   string             `yaml:"version"`
	Pipelines Pipelines          `yaml:"pipelines,omitempty"`
	Actions   map[string]*Action `yaml:"actions,omitempty"`
	WordLists map[string]string  `yaml:"wordlists,omitempty"`
	Tools     []*Tool            `yaml:"tools,omitempty"`
}

// Pipelines represent a collection of pipelines
type Pipelines map[string]*Pipeline

// Pipeline represents the document saving a user's workflow jobs when investigating a target
type Pipeline struct {
	Name        string
	Description string          `yaml:"description,omitempty"`
	Jobs        map[string]*Job `yaml:"jobs,omitempty"`
}

// Validate verifies a new pipeline entry
func (p *Pipeline) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("pipelines must have a name")
	}

	if strings.ContainsAny(p.Name, " ") {
		return fmt.Errorf("pipeline names cannot have spaces")
	}

	if p.Description == "" {
		return fmt.Errorf("pipelines must have a brief description - you'll thank me later")
	}

	return nil
}
