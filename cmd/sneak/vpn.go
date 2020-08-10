package sneak

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var openVPN *vpn.OpenVPN
var home string

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

		// ctx := context.Background()
		// runner := shell.NewRunner()
		err := openVPN.Connect()
		if err != nil {
			gui.ExitWithError(err)
		}

		// vpnPath := fmt.Sprintf("%s/vpn", home)
		// gui.Info("popcorn", "trying to connect...", vpnPath)
		// v := viper.GetViper()
		// color.Green("open vpn config filepath at %s", v.GetString("openvpn_filepath"))

		// gui.Spin.Start()
		// // vpnscript := fmt.Sprintf("sudo %s %s", vpnPath, v.GetString("openvpn_filepath"))
		// // logs := execute(exec.Command("/bin/bash", "-c", "sudo", vpnPath, v.GetString("openvpn_filepath")))
		// logs, err := exec.Command("/bin/bash", "-c", vpnPath, v.GetString("openvpn_filepath")).CombinedOutput()
		// gui.Spin.Stop()
		// if err != nil {
		// 	gui.ExitWithError(string(logs))
		// }

		// color.Green(string(logs))
		// logs, err := runner.BashExec(vpnscript)
		// if err != nil {
		// 	gui.ExitWithError(err)
		// }
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
