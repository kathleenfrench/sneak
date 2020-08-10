package htb

// Recon represents tasks run to do recon on an attack target
type Recon struct {
	Tasks []*Task `yaml:"tasks"`
}

// Task represents a check run against a target and its findings
type Task struct {
	Running    bool    `yaml:"running"`
	Complete   bool    `yaml:"complete"`
	Logs       string  `yaml:"logs"`
	OutputPath string  `yaml:"out_path"`
	Vulns      []*Vuln `yaml:"vulns"`
}

// Vuln represents a vulnerability
type Vuln struct {
	Type   string `yaml:"type"`
	Explit string `yaml:"exploit"`
	Script string `yaml:"script"`
	Notes  string `yaml:"notes"`
}

type task struct {
	name     string
	command  string
	running  bool
	complete bool
	outPath  string
}
