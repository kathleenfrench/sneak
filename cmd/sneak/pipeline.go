package sneak

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/htb"
	"github.com/kathleenfrench/sneak/internal/repository"
	pRepo "github.com/kathleenfrench/sneak/internal/repository/pipeline"
	"github.com/kathleenfrench/sneak/internal/usecase/pipeline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pipelineUsecase    pipeline.Usecase
	pipelineRepository repository.PipelineRepository
	pipelineGUI        *htb.PipelineGUI
)

var pipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"p", "pip", "pipe", "pipelines", "ps"},
	Short:   "pipelines are a collection of actions defined by the user for running various workflows",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		manifestPath := fmt.Sprintf("%s/manifest.yaml", viper.GetString("cfg_dir"))
		pipelineRepository = pRepo.NewPipelineRepository(manifestPath)
		pipelineUsecase = pipeline.NewPipelineUsecase(pipelineRepository)
		pipelineGUI = htb.NewPipelineGUI(pipelineUsecase)
		manifestExists, err := pipelineUsecase.ManifestExists()
		switch {
		case err != nil:
			return err
		case manifestExists:
			return err
		default:
			err = pipelineUsecase.NewManifest()
			if err != nil {
				return err
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var pipelineNewCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"add", "create"},
	Short:   "add a new pipeline to your sneak workflow",
	Run: func(cmd *cobra.Command, args []string) {
		newPipeline, err := pipelineGUI.PromptUserForPipelineData()
		if err != nil {
			gui.ExitWithError(err)
		}

		err = pipelineUsecase.SavePipeline(newPipeline)
		if err != nil {
			gui.ExitWithError(err)
		}

		gui.Info("+1", fmt.Sprintf("%s was added successfully!", newPipeline.Name), newPipeline.Name)
	},
}

var pipelineListCmd = &cobra.Command{
	Use:     "list",
	Short:   "list all of your pipelines",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		pipelines, err := pipelineUsecase.GetAll()
		if err != nil {
			gui.ExitWithError(err)
		}

		if len(pipelines) == 0 {
			gui.Warn("you don't have any pipelines configured yet! run `sneak pipeline new` to get started", nil)
			return
		}

		selection := pipelineGUI.SelectPipelineFromDropdown(pipelines)

		if err = pipelineGUI.SelectPipelineActionsDropdown(selection, pipelines); err != nil {
			gui.ExitWithError(err)
		}
	},
}

var pipelineManifestActionsCmd = &cobra.Command{
	Use:     "actions",
	Aliases: []string{"action", "act"},
	Short:   "define common actions for re-use between multiple pipelines in your pipeline manifest",
	Run: func(cmd *cobra.Command, args []string) {
		color.Red("todo")
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

func init() {
	pipelineCmd.AddCommand(pipelineNewCmd)
	pipelineCmd.AddCommand(pipelineListCmd)
	pipelineCmd.AddCommand(pipelineManifestActionsCmd)
	pipelineCmd.AddCommand(pipelineManifestWordlistsCmd)
}
