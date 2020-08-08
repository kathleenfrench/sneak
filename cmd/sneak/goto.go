package sneak

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var gotoCmd = &cobra.Command{
	Use:   "goto",
	Short: "open your browser to various resources",
	Run: func(cmd *cobra.Command, args []string) {
		color.Red("TODO")
	},
}
