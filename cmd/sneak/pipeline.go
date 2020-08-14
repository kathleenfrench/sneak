package sneak

import (
	"fmt"

	"github.com/fatih/color"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		newPipeline, err := pipelineGUI.PromptUserForPipelineData()
		if err != nil {
			return err
		}

		color.Green("new pipeline: %v", newPipeline)
		return nil
	},
}

func init() {
	pipelineCmd.AddCommand(pipelineNewCmd)
}
