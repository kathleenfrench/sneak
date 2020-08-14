package htb

import (
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
)

// PromptUserForPipelineData prompts the user for val.ues about the pipeline they want to add
func (pg *PipelineGUI) PromptUserForPipelineData() (*entity.Pipeline, error) {
	p := &entity.Pipeline{
		Name:        gui.InputPromptWithResponse("what do you want to call this pipeline?", "", true),
		Description: gui.InputPromptWithResponse("provide a brief description of what this pipeline does for your reference later", "", true),
	}

	err := p.Validate()
	if err != nil {
		return nil, err
	}

	return p, nil
}
