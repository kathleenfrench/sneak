package sneak

import "github.com/spf13/cobra"

var pipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"p", "pip", "pipe", "pipelines", "ps"},
	Short:   "pipelines are a collection of actions defined by the user for running various workflows",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}
