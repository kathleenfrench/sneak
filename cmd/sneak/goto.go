package sneak

import (
	"github.com/kathleenfrench/common/exec"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/helpers"
	"github.com/spf13/cobra"
)

var gotoCmd = &cobra.Command{
	Use:   "goto",
	Short: "open your browser to various resources",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shortcutKeys := helpers.GetKeysFromMap(sneakCfg.WebShortcuts)
		switch len(args) {
		case 0:
			// dropdown
			choice := gui.SelectPromptWithResponse("select a shortcut", shortcutKeys, nil, true)
			err := exec.OpenURL(sneakCfg.WebShortcuts[choice])
			if err != nil {
				gui.ExitWithError(err)
			}

			return
		default:
			// accept one the set keys

		}
	},
}
