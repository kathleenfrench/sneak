package sneak

import (
	"fmt"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/htb"
	"github.com/kathleenfrench/sneak/internal/repository/box"
	boxusecase "github.com/kathleenfrench/sneak/internal/usecase/box"
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
		boxRepo := box.NewBoxRepository(db)
		boxUsecase := boxusecase.NewUsecase(boxRepo)
		nb := htb.NewBoxGUI(boxUsecase)

		box, err := nb.PromptUserForBoxData()
		if err != nil {
			gui.ExitWithError(err)
		}

		err = boxUsecase.Save(box)
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
		br := box.NewBoxRepository(db)
		bu := boxusecase.NewUsecase(br)

		boxes, err := bu.GetAll()
		if err != nil {
			gui.ExitWithError(err)
		}

		if len(boxes) == 0 {
			gui.Warn("you don't have any boxes yet! run `sneak box add` to get started", nil)
			return
		}

		nb := htb.NewBoxGUI(bu)
		selection := nb.SelectBoxFromDropdown(boxes)

		if err = nb.SelectBoxActionsDropdown(selection, boxes); err != nil {
			gui.ExitWithError(err)
		}
	},
}

func init() {
	boxSubCmd.AddCommand(newBoxCmd)
	boxSubCmd.AddCommand(listBoxesCmd)
}
