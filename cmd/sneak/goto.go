package sneak

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var gotoCmd = &cobra.Command{
	Use:   "goto",
	Short: "open your browser to various resources",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shortcuts := sneakCfg.WebShortcuts

		color.Green("shortcuts: %v", shortcuts)

		switch len(args) {
		case 0:
			// dropdown
		default:
			// accept one the set keys

		}
	},
}
