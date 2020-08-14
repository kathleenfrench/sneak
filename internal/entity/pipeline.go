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
	WordLists []string           `yaml:"wordlist_paths,omitempty"`
}

// Pipelines represent a collection of pipelines
type Pipelines map[string]*Pipeline

// Pipeline represents the document saving a user's workflow jobs when investigating a target
type Pipeline struct {
	Name        string
	Description string          `yaml:"description,omitempty"`
	Jobs        map[string]*Job `yaml:"jobs,omitempty"`
}

// Job represents a collection of tasks to run
type Job struct {
	Name  string           `yaml:"name"`
	Tasks map[string]*Task `yaml:"tasks,omitempty"`
	Tools []Tool           `yaml:"tools,omitempty"`
}

// Task represents an individual check
type Task struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	Runner      *Runner `yaml:"run,omitempty"`
	ScriptPath  string  `yaml:"script_path,omitempty"`
	OutputPath  string  `yaml:"output_path,omitempty"`
	Skip        bool    `yaml:"skip,omitempty"`
	Action      *Action `yaml:"action,omitempty"`
	output      string
	complete    bool
	active      bool
}

// Runner represents the actual scripting/tools to use when running the task
type Runner struct {
	Command    string `yaml:"command,omitempty"`
	ScriptPath string `yaml:"script_path,omitempty"`
	OutputPath string `yaml:"output_path,omitempty"`
	output     string
	active     bool
	completed  bool
}

// Tool represents information about an external tool to use/download
type Tool struct {
	Name                 string `yaml:"name"`
	Description          string `yaml:"description"`
	DownloadURL          string `yaml:"download_url,omitempty"`
	ScriptPath           string `yaml:"script_path,omitempty"`
	DownloadIfNotPresent bool   `yaml:"download_if_not_present,omitempty"`
}

// Action represents an independent action that can be taken - useful for common operations like reverse shells, spinning up simple http servers, reverse shells, etc.
type Action struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Runner      Runner `yaml:"run"`
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
