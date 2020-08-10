package sneak

import (
	"fmt"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/pkg/htb"
	"github.com/spf13/cobra"
)

// todo: add dropdown of options - add, list, etc
var boxSubCmd = &cobra.Command{
	Use:     "box",
	Aliases: []string{"boxes", "b"},
	Short:   "do stuff with htb boxes",
}

var newBoxCmd = &cobra.Command{
	Use:     "new",
	Short:   "add a new box",
	Aliases: []string{"add", "a"},
	Run: func(cmd *cobra.Command, args []string) {
		box, err := htb.PromptUserForBoxData()
		if err != nil {
			gui.ExitWithError(err)
		}

		err = htb.SaveBox(db, box)
		if err != nil {
			gui.ExitWithError(err)
		}

		gui.Info("+1", fmt.Sprintf("%s was added successfully!", box.Name), fmt.Sprintf("%s", box.IP))
	},
}

var listBoxesCmd = &cobra.Command{
	Use:     "list",
	Short:   "list all of your boxes",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		boxes, err := htb.GetAllBoxes(db)
		if err != nil {
			gui.ExitWithError(err)
		}

		if len(boxes) == 0 {
			gui.Warn("you don't have any boxes yet! run `sneak box add` to get started", nil)
			return
		}

		selection := htb.SelectBoxFromDropdown(boxes)

		if err = htb.SelectBoxActionsDropdown(db, selection, boxes); err != nil {
			gui.ExitWithError(err)
		}
	},
}

func init() {
	boxSubCmd.AddCommand(newBoxCmd)
	boxSubCmd.AddCommand(listBoxesCmd)
}
