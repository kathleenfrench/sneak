package entity

// Action represents an independent action that can be taken - useful for common operations like reverse shells, spinning up simple http servers, reverse shells, etc.
type Action struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Runner      *Runner `yaml:"run"`
}

// Dependency represents what dependency or dependencies are needed to run a given action
type Dependency struct {
	Name           string `yaml:"name"`
	DownloadLink   string `yaml:"download_url"`
	DownloadScript string `yaml:"download_script"`
}
