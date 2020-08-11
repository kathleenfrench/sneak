package pipeline

// Workflow represents the document saving a user's preset and custom pipelines
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
	Name       string `yaml:"name"`
	Runner     Runner `yaml:"command,omitempty"`
	ScriptPath string `yaml:"script_path,omitempty"`
	Output     string `yaml:"output,omitempty"`
	OutputPath string `yaml:"output_path,omitempty"`
	Complete   bool   `yaml:"complete,omitempty"`
	InProgress bool   `yaml:"in_progress,omitempty"`
	Skip       bool   `yaml:"skip,omitempty"`
}

// Runner represents the actual scripting/tools to use when running the task
type Runner struct {
	Command    string `yaml:"command,omitempty"`
	ScriptPath string `yaml:"script_path,omitempty"`
	Active     bool   `yaml:"active,omitempty"`
	Complete   bool   `yaml:"complete,omitempty"`
	Output     string `yaml:"output,omitempty"`
	OutputPath string `yaml:"output_path,omitempty"`
}

// Tool represents information about an external tool to use/download
type Tool struct {
	DownloadURL string `yaml:"download_url,omitempty"`
	ScriptPath  string `yaml:"script_path,omitempty"`
}
