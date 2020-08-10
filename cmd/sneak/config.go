package sneak

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/exec"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"configs", "c"},
	Short:   "view and/or modify your sneak config values",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rootCmd.PersistentPreRun(cmd, args)
	},
}

var configListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "show", "view"},
	Short:   "view your current sneak configs",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.BashExec(fmt.Sprintf("cat %s", viper.ConfigFileUsed()))
		if err != nil {
			gui.ExitWithError(err)
		}

		color.HiBlue(out)
	},
}

var configGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "get a specific config value",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.GetViper()
		if v.IsSet(args[0]) {
			color.Green("%v", v.Get(args[0]))
		} else {
			color.Red("no value set for key %s!", args[0])
		}
	},
}

var configUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "change"},
	Args:    cobra.MaximumNArgs(0),
	Short:   "modify your existing sneak configs",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.UpdateSettingsPrompt(viper.AllSettings())
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var configDelCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "delete non-required config keys and values",
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configUpdateCmd)
	configCmd.AddCommand(configDelCmd)
	configCmd.AddCommand(configGetCmd)
}
