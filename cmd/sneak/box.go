package sneak

import (
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/htb"
	"github.com/kathleenfrench/sneak/internal/repository"
	"github.com/kathleenfrench/sneak/internal/repository/box"
	boxusecase "github.com/kathleenfrench/sneak/internal/usecase/box"
	"github.com/spf13/cobra"
)

var (
	boxUsecase    boxusecase.Usecase
	boxRepository repository.BoxRepository
	boxGUI        *htb.BoxGUI
)

// todo: add dropdown of options - add, list, etc
var boxSubCmd = &cobra.Command{
	Use:     "box",
	Aliases: []string{"boxes", "b"},
	Short:   "do stuff with htb boxes",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		boxRepository = box.NewBoxRepository(db)
		boxUsecase = boxusecase.NewUsecase(boxRepository)
		boxGUI = htb.NewBoxGUI(boxUsecase)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := boxGUI.DefaultDropdownHandler()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var newBoxCmd = &cobra.Command{
	Use:     "new",
	Short:   "add a new box",
	Aliases: []string{"add", "a"},
	Run: func(cmd *cobra.Command, args []string) {
		err := boxGUI.AddBox()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var listBoxesCmd = &cobra.Command{
	Use:     "list",
	Short:   "list all of your boxes",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		err := boxGUI.ListBoxes()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

func init() {
	boxSubCmd.AddCommand(newBoxCmd)
	boxSubCmd.AddCommand(listBoxesCmd)
}
