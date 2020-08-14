package sneak

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/helpers"
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

		err := openVPN.Connect()
		if err != nil {
			gui.ExitWithError(err)
		}

		color.Green("success! you can always run `sneak vpn test` to verify your connection to the HTB lab IP set in your configs")
	},
}

func execute(cmd *exec.Cmd) string {
	var out string

	if cmd == nil {
		return out
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		out = color.HiRedString(fmt.Sprintf("[ERROR]: %s", err))
		return out
	}

	err = cmd.Start()
	if err != nil {
		out = color.HiRedString(fmt.Sprintf("[ERROR]: %s", err))
		return out
	}

	res, err := ioutil.ReadAll(stdout)
	if err != nil {
		out = color.HiRedString(fmt.Sprintf("[ERROR]: %s", err))
		return out
	}

	out += string(res)

	err = cmd.Wait()
	if err != nil {
		out = color.HiRedString(fmt.Sprintf("[ERROR]: %s", err))
		return out
	}

	return strings.TrimSpace(out)
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
	Aliases: []string{"ping", "check", "status"},
	Run: func(cmd *cobra.Command, args []string) {
		if !openVPN.AlreadySetup() {
			gui.ExitWithError("you have not setup openvpn yet with sneak - run 'sneak vpn setup'")
		}

		if !viper.IsSet("htb_network_ip") {
			gui.ExitWithError("your HTB network IP has not been set in your configs - please run `sneak config update`, select it from the dropdown, and follow the onscreen prompts")
		}

		gui.Info("popcorn", fmt.Sprintf("\033[H\033[2J\nchecking..."), viper.GetString("htb_network_ip"))
		gui.Spin.Start()
		err := helpers.SudoPing(viper.GetString("htb_network_ip"))
		gui.Spin.Stop()
		if err != nil {
			gui.Warn("your connection could not be established", viper.Get("htb_network_ip"))
			gui.ExitWithError(err)
		}

		color.Green("you're connected!")
	},
}

func init() {
	vpnCmd.AddCommand(vpnSetupCmd)
	vpnCmd.AddCommand(vpnUpdateCmd)
	vpnCmd.AddCommand(vpnTestCmd)
	vpnCmd.AddCommand(vpnConnectCmd)
}
