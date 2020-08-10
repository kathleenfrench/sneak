package sneak

import (
	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var openVPN *vpn.OpenVPN

var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "connect to the openvpn client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		openVPN, err = vpn.NewOpenVPNClient()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var vpnSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup sneak to work with openvpn",
	Run: func(cmd *cobra.Command, args []string) {
		if openVPN.AlreadySetup() {
			color.Green("your openvpn configs have alread been setup!")
			return
		}

		err := openVPN.Setup(viper.GetString("default_editor"))
		if err != nil {
			gui.ExitWithError(err)
		}

		gui.Info("+1", "your openvpn config file has been set!", openVPN.Filepath)
	},
}

var vpnConnectCmd = &cobra.Command{
	Use:     "connect",
	Short:   "connect to the htb servers",
	Aliases: []string{"start"},
	Run: func(cmd *cobra.Command, args []string) {
		if !openVPN.AlreadySetup() {
			gui.ExitWithError("you have not setup openvpn yet with sneak - run 'sneak vpn setup'")
		}

	},
}

var vpnUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "change your openvpn config file",
	Aliases: []string{"u", "change", "edit"},
	Run: func(cmd *cobra.Command, args []string) {
		if !openVPN.AlreadySetup() {
			gui.ExitWithError("you have not setup openvpn yet with sneak - run 'sneak vpn setup'")
		}

		changeContentsFile := gui.ConfirmPrompt("do you want to change the contents of the .ovpn file?", "", false, true)

		switch changeContentsFile {
		case true:
			err := openVPN.Setup(viper.GetString("default_editor"))
			if err != nil {
				gui.ExitWithError(err)
			}

			gui.Info("+1", "your openvpn config file has been updated!", openVPN.Filepath)
		default:
			color.HiBlue("leaving it as is for now, then...")
		}
	},
}

var vpnTestCmd = &cobra.Command{
	Use:     "test",
	Short:   "test your vpn connection",
	Aliases: []string{"ping", "check"},
	Run: func(cmd *cobra.Command, args []string) {
		if !openVPN.AlreadySetup() {
			gui.ExitWithError("you have not setup openvpn yet with sneak - run 'sneak vpn setup'")
		}
	},
}

func init() {
	vpnCmd.AddCommand(vpnSetupCmd)
	vpnCmd.AddCommand(vpnUpdateCmd)
	vpnCmd.AddCommand(vpnTestCmd)
	vpnCmd.AddCommand(vpnConnectCmd)
}
