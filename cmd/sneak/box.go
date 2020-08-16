package sneak

import (
	"fmt"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/htb"
	"github.com/kathleenfrench/sneak/internal/repository"
	"github.com/kathleenfrench/sneak/internal/repository/box"
	pRepo "github.com/kathleenfrench/sneak/internal/repository/pipeline"
	boxusecase "github.com/kathleenfrench/sneak/internal/usecase/box"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		manifestPath := fmt.Sprintf("%s/manifest.yaml", viper.GetString("cfg_dir"))
		boxRepository = box.NewBoxRepository(db)
		boxUsecase = boxusecase.NewUsecase(boxRepository)
		pipelineRepository = pRepo.NewPipelineRepository(manifestPath)
		pipelineUsecase = pipeline.NewPipelineUsecase(pipelineRepository)
		boxGUI = htb.NewBoxGUI(boxUsecase, pipelineUsecase)
		manifestExists, err := pipelineUsecase.ManifestExists()
		switch {
		case err != nil:
			gui.ExitWithError(err)
		case manifestExists:
			return
		default:
			err = pipelineUsecase.NewManifest()
			if err != nil {
				gui.ExitWithError(err)
			}
		}
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
		box, err := boxGUI.PromptUserForBoxData()
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
		boxes, err := boxUsecase.GetAll()
		if err != nil {
			gui.ExitWithError(err)
		}

		if len(boxes) == 0 {
			gui.Warn("you don't have any boxes yet! run `sneak box new` to get started", nil)
			return
		}

		selection := boxGUI.SelectBoxFromDropdown(boxes)

		if err = boxGUI.SelectBoxActionsDropdown(selection, boxes); err != nil {
			gui.ExitWithError(err)
		}
	},
}

func init() {
	boxSubCmd.AddCommand(newBoxCmd)
	boxSubCmd.AddCommand(listBoxesCmd)
}
