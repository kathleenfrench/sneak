package sneak

import (
	"os"

	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version is a value injected at compile time for the current version of sneak
var Version = "master"

// local
var sneakCfg config.Settings

var rootCmd = &cobra.Command{
	Use:     "sneak",
	Aliases: []string{"snk"},
	Short:   "a tool for common actions when pentesting/playing CTFs",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		s, err := config.Parse(viper.GetViper())
		if err != nil {
			os.Exit(0)
		}

		sneakCfg = s
	},
}
