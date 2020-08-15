package entity

// Job represents a collection of tasks to run
type Job struct {
	Name        string
	Description string             `yaml:"description,omitempty"`
	OneOffs     map[string]*Runner `yaml:"oneoffs,omitempty"`
	Actions     []string           `yaml:"actions,omitempty"` // where string = action name
	Disabled    bool               `yaml:"disabled,omitempty"`
}
