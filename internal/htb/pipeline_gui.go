package htb

import "github.com/kathleenfrench/sneak/internal/usecase/pipeline"

// PipelineGUI manages methods and properties of the terminal GUI re: pipelines
type PipelineGUI struct {
	usecase pipeline.Usecase
}

// NewPipelineGUI instantiates a new pipeline gui struct
func NewPipelineGUI(u pipeline.Usecase) *PipelineGUI {
	return &PipelineGUI{
		usecase: u,
	}
}
