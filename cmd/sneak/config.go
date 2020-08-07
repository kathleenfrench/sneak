package sneak

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/exec"
	"github.com/kathleenfrench/common/finder"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"configs", "c"},
	Short:   "view and/or modify your sneak config values",
	Args:    cobra.MinimumNArgs(1),
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
		key := args[0]
		val := args[1]

		keys := viper.AllKeys()
		// filter out: useviper, cfg_dir
		settings := viper.AllSettings()

		var toChange interface{}
		if finder.FoundStringInSlice(keys, key) {
			toChange = settings[key]
		}

		color.HiBlue("will be changing %v to %s", toChange, val)

		// switch len(args) {
		// case 1:
		// 	// just receive the key, verify that they key exists and prompt for what value to change
		// 	key := args[0]
		// 	keys := v.AllKeys()
		// 	if finder.FoundStringInSlice(keys, key) {
		// 		gui.Info("eyes", fmt.Sprintf("%s is an existing key...", key), v.Get(key))
		// 		// use the currently set value as the default
		// 	} else {

		// 	}
		// case 2:
		// 	// receive the key and the value, set the key and value
		// default:
		// 	gui.ExitWithError("invalid input")
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

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUpdateCmd)
	configCmd.AddCommand(configDelCmd)
}
