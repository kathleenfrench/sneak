package sneak

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/htb"
	"github.com/kathleenfrench/sneak/internal/repository"
	pRepo "github.com/kathleenfrench/sneak/internal/repository/pipeline"
	"github.com/kathleenfrench/sneak/internal/usecase/action"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pipelineUsecase    pipeline.Usecase
	pipelineRepository repository.PipelineRepository
	pipelineGUI        htb.PipelineGUI
)

var (
	actionUsecase action.Usecase
	actionGUI     *htb.ActionsGUI
)

var pipelineCmd = &cobra.Command{
	Use:     "pipelines",
	Aliases: []string{"p", "pip", "pipe", "pipeline", "ps"},
	Short:   "pipelines are a collection of actions defined by the user for running various workflows",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		manifestPath := fmt.Sprintf("%s/manifest.yaml", viper.GetString("cfg_dir"))
		pipelineRepository = pRepo.NewPipelineRepository(manifestPath)
		pipelineUsecase = pipeline.NewPipelineUsecase(pipelineRepository)
		pipelineGUI = htb.NewPipelineGUI(pipelineUsecase)
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
		err := pipelineGUI.DefaultPipelineDropdownHandler()
		strErr := err.Error()
		switch strErr {
		case "actions":
			pipelineManifestActionsCmd.Run(nil, nil)
			return
		}

		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var pipelineNewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"add", "create"},
	Short:   "add a new pipeline to your sneak workflow",
	Run: func(cmd *cobra.Command, args []string) {
		err := pipelineGUI.AddPipeline()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var pipelineListCmd = &cobra.Command{
	Use:     "list",
	Short:   "list all of your pipelines",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		err := pipelineGUI.ListPipelines()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var pipelineManifestWordlistsCmd = &cobra.Command{
	Use:     "wordlists",
	Aliases: []string{"word", "w", "wl", "wordlist"},
	Short:   "add wordlists for re-use between multiple pipelines in your pipeline manifest",
	Run: func(cmd *cobra.Command, args []string) {
		color.Red("todo")
	},
}

var pipelineManifestActionsCmd = &cobra.Command{
	Use:     "actions",
	Aliases: []string{"action", "act"},
	Short:   "define common actions for re-use between multiple pipelines in your pipeline manifest",
	Run: func(cmd *cobra.Command, args []string) {
		err := actionGUI.HandleActionsDropdown()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

func init() {
	pipelineCmd.AddCommand(pipelineNewCmd)
	pipelineCmd.AddCommand(pipelineListCmd)
	pipelineCmd.AddCommand(pipelineManifestWordlistsCmd)
	pipelineCmd.AddCommand(pipelineManifestActionsCmd)
}
