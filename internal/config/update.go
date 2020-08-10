package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/exec"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/helpers"
	"github.com/spf13/viper"
)

// UpdateSettings checks for pls config values that have already been set and ensures they're preserved when updating configs
func (s *Settings) UpdateSettings() error {
	cfgFile := s.viper.ConfigFileUsed()

	if s.HTBUsername != "" {
		s.viper.Set("htb_username", strings.TrimSpace(s.HTBUsername))
	}

	if s.OpenVPNFilepath != "" {
		s.viper.Set("openvpn_filepath", s.OpenVPNFilepath)
	}

	if s.DefaultEditor != "" {
		s.viper.Set("default_editor", s.DefaultEditor)
	}

	if s.WebShortcuts != nil {
		s.viper.Set("webshort", s.WebShortcuts)
	}

	if s.HTBNetworkIP != "" {
		s.viper.Set("htb_network_ip", s.HTBNetworkIP)
	}

	s.viper.MergeInConfig()
	s.viper.SetConfigFile(cfgFile)
	s.viper.SetConfigType(filepath.Ext(cfgFile))
	err := s.viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func filterForChangeableConfigs(keys []string) []string {
	eligible := []string{}
	for _, k := range keys {
		switch k {
		case "useviper", "cfg_dir":
			break
		default:
			if strings.Contains(k, "webshort.") {
				break
			}

			eligible = append(eligible, k)
		}
	}

	eligible = append(eligible, "webshort")
	return eligible
}

// UpdateSettingsPrompt creates a dropdown in the terminal UI for the user to select which config value to change
func UpdateSettingsPrompt(viperSettings map[string]interface{}) error {
	var changedValue string

	v := viper.GetViper()
	keys := filterForChangeableConfigs(viper.AllKeys())
	choice := gui.SelectPromptWithResponse("which config value do you want to change?", keys, nil, false)

	switch choice {
	case "default_editor":
		changedValue = gui.GetUsersPreferredEditor("", true)
		v.Set(choice, changedValue)
	case "htb_network_ip":
		fmt.Println(htbNetworkIPHelpText)
		changedValue = gui.InputPromptWithResponse("what is your HTB Lab Network IPv4?", "", true)
		v.Set(choice, changedValue)
	case "webshort":
		shorts := make(map[string]string)
		err := viper.UnmarshalKey("webshort", &shorts)
		if err != nil {
			gui.ExitWithError(err)
		}

		editExisting := gui.ConfirmPrompt("do you want to modify an existing url?", "", false, true)
		if editExisting {
			shortKeys := helpers.GetKeysFromMap(shorts)
			editWhich := gui.SelectPromptWithResponse("which do you want to change?", shortKeys, nil, false)
			changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", editWhich), "", true)
			color.Green("setting %s", fmt.Sprintf("webshort.%s", editWhich))
			v.Set(fmt.Sprintf("webshort.%s", editWhich), changedValue)
		} else {
			target, url := addNewWebShortcut()
			v.Set(target, url)
			changedValue = url
		}
	default:
		switch strings.Contains(choice, ".") {
		case true:
			color.Red("TODO")
		default:
			changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice), fmt.Sprintf("%v", v.Get(choice)), true)
			v.Set(choice, changedValue)
		}
	}

	parsed, err := Parse(v)
	if err != nil {
		gui.ExitWithError(err)
	}

	parsed.UpdateSettings()
	color.HiGreen("successfully updated %s to equal %s", choice, changedValue)
	return nil
}

func addNewWebShortcut() (string, string) {
	target := gui.InputPromptWithResponse("what do you want to name the shortcut?", "", false)

	url := gui.InputPromptWithResponse(fmt.Sprintf("what is the shortcut url you want to set for %s?", target), "", false)

	if target == "" || url == "" {
		gui.ExitWithError("missing required values")
	}

	return fmt.Sprintf("webshort.%s", target), url
}

func verify(filepath string, mountedData bool) error {
	preCheck, _ := Parse(viper.GetViper())

	if preCheck.WebShortcuts == nil {
		viper.Set("webshort", defaultShortcuts)
	}

	if preCheck.HTBNetworkIP == "" {
		fmt.Println(htbNetworkIPHelpText)
		htbIP := gui.InputPromptWithResponse("what is your HTB Lab Network IPv4?", "", true)
		viper.Set("htb_network_ip", htbIP)
	}

	if mountedData {
		if whoami, err := exec.BashExec("whoami"); err != nil {
			gui.Warn("sneak had an issue updating your ovpn_filepath config from what's coming from your mounted data - make sure to run `sneak config update` and change it so you can connect to the vpn", err)
		} else {
			newOvpnPath := fmt.Sprintf("/home/%s/.sneak/%s.ovpn", whoami, whoami)
			viper.Set("openvpn_filepath", newOvpnPath)
		}
	}

	return viper.WriteConfigAs(filepath)
}

var htbNetworkIPHelpText = color.YellowString(`
your HTB lab access network information can be found at:
https://www.hackthebox.eu/home/htb/access

**important**: whenever you change servers you will
need to re-download your connection pack and update
the .ovpn file in sneak as well as the 'HTB Network IPv4'
value used by sneak to test your connection.
`)
