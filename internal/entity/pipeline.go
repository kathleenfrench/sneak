package entity

// Pipelines represent a collection of pipelines
type Pipelines map[string]Pipeline

// Pipeline represents the document saving a user's workflow jobs when investigating a target
type Pipeline struct {
	Version   int               `yaml:"version"`
	Jobs      map[string]Job    `yaml:"jobs,omitempty"`
	WordLists []string          `yaml:"wordlist_paths,omitempty"`
	Actions   map[string]Action `yaml:"actions,omitempty"`
}

// Job represents a collection of tasks to run
type Job struct {
	Name  string          `yaml:"name"`
	Tasks map[string]Task `yaml:"tasks,omitempty"`
	Tools []Tool          `yaml:"tools,omitempty"`
}

// Task represents an individual check
type Task struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Runner      Runner `yaml:"run,omitempty"`
	ScriptPath  string `yaml:"script_path,omitempty"`
	OutputPath  string `yaml:"output_path,omitempty"`
	Skip        bool   `yaml:"skip,omitempty"`
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
