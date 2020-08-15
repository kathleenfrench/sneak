package entity

// Runner represents the actual scripting/tools to use when running the task
type Runner struct {
	Command      string `yaml:"command,omitempty"`
	ScriptPath   string `yaml:"script_path,omitempty"`
	DontSaveLogs bool   `yaml:"disable_logsave,omitempty"`
	logs         string
	active       bool
	completed    bool
}
