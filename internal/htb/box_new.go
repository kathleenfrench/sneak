package htb

import (
	"fmt"
	"time"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/txn2/txeh"
)

// PromptUserForBoxData prompts the user for values about the htb machine they want to add
func (bg *BoxGUI) PromptUserForBoxData() (entity.Box, error) {
	box := entity.Box{
		Name:        gui.InputPromptWithResponse("what is the name of the box?", "", true),
		IP:          gui.InputPromptWithResponse("what is its IP?", "", true),
		Description: gui.InputPromptWithResponse("provide a short description of the box if you want", "", true),
		Completed:   false,
		Active:      false,
		Notes:       "",
		OS:          gui.SelectPromptWithResponse("what is the OS?", osOptions, nil, true),
		Difficulty:  gui.SelectPromptWithResponse("what is its difficulty?", difficulties, nil, true),
		Up:          false,
		Flags: entity.Flags{
			Root: "",
			User: "",
		},
		Created:     time.Now(),
		LastUpdated: time.Now(),
	}

	if err := box.Validate(); err != nil {
		return box, err
	}

	box.Hostname = fmt.Sprintf("%s.htb", box.Name)
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return box, err
	}

	hosts.AddHost(box.IP, box.Hostname)
	return box, nil
}
