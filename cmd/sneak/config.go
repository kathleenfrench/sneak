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
	Args:    cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.BashExec(fmt.Sprintf("cat %s", viper.ConfigFileUsed()))
		if err != nil {
			gui.ExitWithError(err)
		}

		color.HiBlue(out)
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

var configSetCmd = &cobra.Command{
	Use:   "set",
	Args:  cobra.ExactArgs(2),
	Short: "skip the update GUI and either change a specific config value by config key or add a custom config",
	Run: func(cmd *cobra.Command, args []string) {
		// key := args[0]
		// val := args[1]

		// switch viper.IsSet(key) {
		// case true:
		// 	// prompt use what they want to change it to
		// 	gui.Info("light_bulb", fmt.Sprintf("%s has a current value of %v - changing it to %s", key, viper.Get(key), val), nil)
		// 	viper.Set(key, val)
		// 	sneakCfg.UpdateSettings()
		// default:
		// 	correct := gui.ConfirmPrompt(fmt.Sprintf("%s is not an existing key - are you meaning to add a new one?", key), "", true, true)
		// 	switch correct {
		// 	case true:
		// 		viper.Set(key, val)
		// 		err := sneakCfg.UpdateSettings()
		// 		if err != nil {
		// 			gui.ExitWithError(err)
		// 		}
		// 	default:
		// 	}
		// }
	},
}

var configDelCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "delete non-required config keys and values",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
