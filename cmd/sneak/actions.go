package sneak

import (
	"fmt"

	"github.com/fatih/color"
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
)

var pipelineManifestActionsCmd = &cobra.Command{
	Use:     "actions",
	Aliases: []string{"action", "act"},
	Short:   "define common actions for re-use between multiple pipelines in your pipeline manifest",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		manifestPath := fmt.Sprintf("%s/manifest.yaml", viper.GetString("cfg_dir"))
		pipelineRepository = pRepo.NewPipelineRepository(manifestPath)
		pipelineUsecase = pipeline.NewPipelineUsecase(pipelineRepository)
		pipelineGUI = htb.NewPipelineGUI(pipelineUsecase)
		actionUsecase = action.NewActionUsecase(pipelineUsecase)
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
		color.Red("todo")
	},
}
