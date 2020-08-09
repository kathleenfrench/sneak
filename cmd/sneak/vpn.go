package sneak

import (
	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/helpers"
	"github.com/spf13/cobra"
)

var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "connect to the openvpn client",
	Run: func(cmd *cobra.Command, args []string) {
		// sudo check
		openvpn, err := helpers.NewOpenVPNClient()
		if err != nil {
			gui.ExitWithError(err)
		}

		color.Green("local network: %s", openvpn.LocalNetwork)
		color.Green("filepath: %s", openvpn.Filepath)
	},
}
