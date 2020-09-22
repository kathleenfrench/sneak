package sneak

import (
	"fmt"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/htb"
	pRepo "github.com/kathleenfrench/sneak/internal/repository/pipeline"
	"github.com/kathleenfrench/sneak/internal/usecase/action"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	actionUsecase action.Usecase
	actionGUI     *htb.ActionsGUI
)

var pipelineManifestActionsCmd = &cobra.Command{
	Use:     "actions",
	Aliases: []string{"action", "act"},
	Short:   "define common actions for re-use between multiple pipelines in your pipeline manifest",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		manifestPath := fmt.Sprintf("%s/manifest.yaml", viper.GetString("cfg_dir"))
		pipelineRepository = pRepo.NewPipelineRepository(manifestPath)
		pipelineUsecase = pipeline.NewPipelineUsecase(pipelineRepository)
		actionUsecase = action.NewActionUsecase(pipelineUsecase)
		actionGUI = htb.NewActionsGUI(actionUsecase)
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
		err := actionGUI.HandleActionsDropdown()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}
