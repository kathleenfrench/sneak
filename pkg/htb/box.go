package htb

// Box represents a machine
type Box struct {
	Name       string `yaml:"name"`
	IP         string `yaml:"ip"`
	Hostname   string `yaml:"hostname"`
	Difficulty string `yaml:"difficulty"` // easy, medium, hard, insane
	Flags      *Flags `yaml:"flags"`
	Completed  bool   `yaml:"completed"`
	Active     bool   `yaml:"active"`
	NotesPath  string `yaml:"notes_path"`
}

// Flags represent the flags to find on a machine
type Flags struct {
	Root string `yaml:"root"`
	User string `yaml:"user"`
}
