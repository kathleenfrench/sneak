package config

import (
	"fmt"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const configPathRelativeToHome = ".config/sneak"
const configFileName = "config"
const configFileType = "yaml"

func constructConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		gui.ExitWithError(err)
	}

	return fmt.Sprintf("%s/%s", home, configPathRelativeToHome)
}

// Initialize creates the directory and/or file with defaults for the application's configuration settings
func Initialize() {
	// set fs properties
	viper.AddConfigPath(constructConfigPath())
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.SetConfigFile(fmt.Sprintf("%s/%s.%s", constructConfigPath(), configFileName, configFileType))

	// check for whether the directory and config file already exist
	err := utils.CreateDir(constructConfigPath())
	if err != nil {
		gui.ExitWithError(err)
	}

	err = utils.CreateFile(viper.ConfigFileUsed())
	if err != nil {
		gui.ExitWithError(err)
	}

	viper.AutomaticEnv()

	_ = viper.SafeWriteConfig()
	err = viper.ReadInConfig()
	if err != nil {
		err = viper.WriteConfig()
		if err != nil {
			gui.ExitWithError(err)
		}
	}

	// set defaults
	unsetValuesFound := checkForUnsetRequiredDefaults()
	if unsetValuesFound {
		err = viper.WriteConfig()
		if err != nil {
			gui.ExitWithError(err)
		}
	}

	viper.WatchConfig()

}

func unset(val interface{}) bool {
	if val == nil || val == "" {
		return true
	}

	return false
}

func checkForUnsetRequiredDefaults() bool {
	var unsetFound bool

	if unset(viper.Get("htb_username")) {
		unsetFound = true
		viper.Set("htb_username", "nightwd60")
	}

	if unset(viper.Get("box_ips")) {
		unsetFound = true
		viper.Set("box_ips", make(map[string]string))
	}

	if unset(viper.Get("openvpn_filepath")) {
		unsetFound = true
		viper.Set("openvpn_filepath", fmt.Sprintf("%s/%s.ovpn", constructConfigPath(), viper.Get("htb_username")))
	}

	if unset(viper.Get("default_editor")) {
		unsetFound = true
		preferredEditor := gui.GetUsersPreferredEditor("", true)
		viper.Set("default_editor", preferredEditor)
	}

	return unsetFound
}