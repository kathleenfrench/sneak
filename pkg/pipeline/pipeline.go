package pipeline

// Workflow represents the document saving a user's workflow pipelines when investigating a target
type Workflow struct {
	Version   int                 `yaml:"version"`
	Pipelines map[string]Pipeline `yaml:"pipelines,omitempty"`
	WordLists []string            `yaml:"wordlist_paths,omitempty"`
}

// Pipeline represents a collection of tasks to run
type Pipeline struct {
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
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	DownloadURL string `yaml:"download_url,omitempty"`
	ScriptPath  string `yaml:"script_path,omitempty"`
}

// Action represents an independent action that can be taken - useful for common operations like reverse shells, spinning up simple http servers, reverse shells, etc.
type Action struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Runner      Runner `yaml:"run"`
}
