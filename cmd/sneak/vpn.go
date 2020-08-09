package sneak

import (
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var openVPN *helpers.OpenVPN

var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "connect to the openvpn client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		color.Green("running...")
		var err error
		openVPN, err = helpers.NewOpenVPNClient()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var vpnSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup sneak to work with openvpn",
	Run: func(cmd *cobra.Command, args []string) {
		if fs.FileExists(openVPN.Filepath) {
			color.Green("your openvpn configs have alread been setup!")
			return
		}

		err := fs.CreateFile(openVPN.Filepath)
		if err != nil {
			gui.ExitWithError(err)
		}

		vpnCfgs := gui.TextEditorInputAndSave("copy in your openvpn file from HTB and save", "", viper.GetString("default_editor"))
		err = ioutil.WriteFile(openVPN.Filepath, []byte(vpnCfgs), 0644)
		if err != nil {
			gui.ExitWithError(err)
		}

		gui.Info("+1", "your openvpn config file has been set!", openVPN.Filepath)
	},
}

func init() {
	vpnCmd.AddCommand(vpnSetupCmd)
}
