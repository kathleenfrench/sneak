package entity

// Tool represents information about an external tool to use/download
type Tool struct {
	Name                 string `yaml:"name"`
	Description          string `yaml:"description"`
	DownloadURL          string `yaml:"download_url,omitempty"`
	ScriptPath           string `yaml:"script_path,omitempty"`
	DownloadIfNotPresent bool   `yaml:"download_if_not_present,omitempty"`
}
