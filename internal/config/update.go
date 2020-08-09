package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
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
	case "webshort":
		shorts := make(map[string]string)
		err := viper.UnmarshalKey("webshort", &shorts)
		if err != nil {
			gui.ExitWithError(err)
		}

		editExisting := gui.ConfirmPrompt("do you want to modify an existing url?", "", false, true)
		if editExisting {
			shortKeys := getKeysFromMap(shorts)
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

func getKeysFromMap(m map[string]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func verify(filepath string) error {
	preCheck, _ := Parse(viper.GetViper())

	if preCheck.WebShortcuts == nil {
		viper.Set("webshort", defaultShortcuts)
	}

	return viper.WriteConfigAs(filepath)
}
