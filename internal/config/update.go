package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
	"github.com/kathleenfrench/common/extract"
	"github.com/kathleenfrench/common/gui"
	"github.com/spf13/viper"
)

// UpdateSettings checks for pls config values that have already been set and ensures they're preserved when updating configs
func (s *Settings) UpdateSettings() error {
	cfgFile := s.viper.ConfigFileUsed()

	if s.HTBUsername != "" {
		s.viper.Set("htb_username", strings.TrimSpace(s.HTBUsername))
	}

	if s.BoxIPs != nil {
		s.viper.Set("box_ips", s.BoxIPs)
	}

	if s.OpenVPNFilepath != "" {
		s.viper.Set("openvpn_filepath", s.OpenVPNFilepath)
	}

	if s.DefaultEditor != "" {
		s.viper.Set("default_editor", s.DefaultEditor)
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

// UpdateSettingsPrompt creates a dropdown in the terminal UI for the user to select which config value to change
func UpdateSettingsPrompt(viperSettings map[string]interface{}) error {
	var changedValue string

	v := viper.GetViper()
	keys := viper.AllKeys()
	choice := gui.SelectPromptWithResponse("which config value do you want to change?", keys, nil, false)

	configType := reflect.TypeOf(v.Get(choice))

	switch configType.Kind() {
	case reflect.String:
		switch choice {
		case "default_editor":
			changedValue = gui.GetUsersPreferredEditor("", true)
		default:
			changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice), fmt.Sprintf("%v", v.Get(choice)), true)
		}
	case reflect.Int, reflect.Bool:
		changedValue = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", choice), fmt.Sprintf("%v", v.Get(choice)), true)
	default:
		// prompt for changing one of the existing settings as a map/map interface, etc.
		switch choice {
		case "box_ips":
			ips := make(map[string]string)
			err := viper.UnmarshalKey(choice, &ips)
			if err != nil {
				return err
			}

			var newIP string
			editExisting := gui.ConfirmPrompt("do you want to modify an existing IP?", "", false, true)
			if editExisting {
				shortKeys := extract.KeysFromMapString(ips)
				editWhich := gui.SelectPromptWithResponse("which do you want to change?", shortKeys, nil, false)
				newIP = gui.InputPromptWithResponse(fmt.Sprintf("what do you want to change %s to?", editWhich), "", true)
				choice = fmt.Sprintf("box_ips.%s", editWhich)
			} else {
				// add new box IP
				boxName := gui.InputPromptWithResponse("what is the name of the new box?", "", false)
				ip := gui.InputPromptWithResponse(fmt.Sprintf("what is the IP for %s?", boxName), "", false)
				choice = fmt.Sprintf("box_ips.%s", boxName)
				newIP = ip
			}

			if !govalidator.IsIP(newIP) {
				badIPErr := fmt.Sprintf("%s is not a valid IP!", newIP)
				gui.ExitWithError(errors.New(badIPErr))
			}

			changedValue = newIP
		}
	}

	v.Set(choice, changedValue)

	parsed, err := Parse(v)
	if err != nil {
		gui.ExitWithError(err)
	}

	parsed.UpdateSettings()
	color.HiGreen("successfully updated %s to equal %s", choice, changedValue)
	return nil
}
