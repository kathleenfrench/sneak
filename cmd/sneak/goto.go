package sneak

import (
	"fmt"

	"github.com/kathleenfrench/common/exec"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/pkg/utils"
	"github.com/spf13/cobra"
)

var gotoCmd = &cobra.Command{
	Use:   "goto",
	Short: "open your browser to various resources",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shortcutKeys := utils.GetKeysFromMap(sneakCfg.WebShortcuts)
		switch len(args) {
		case 0:
			choice := gui.SelectPromptWithResponse("select a shortcut", shortcutKeys, nil, true)
			err := exec.OpenURL(sneakCfg.WebShortcuts[choice])
			if err != nil {
				gui.ExitWithError(fmt.Errorf("could not open %s - if you are running sneak in docker, this subcommand will not work - %w", sneakCfg.WebShortcuts[choice], err))
			}

			return
		default:
			// accept one the set keys

		}
	},
}
