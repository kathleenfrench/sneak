package entity

// Action represents an independent action that can be taken - useful for common operations like reverse shells, spinning up simple http servers, reverse shells, etc.
type Action struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Runner      *Runner `yaml:"run"`
}
