package config

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/common/gui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var configName = ".sneak"
var configType = "yaml"

var cfg *Settings

// InitConfig initializes viper with sneak defaults
func InitConfig() {
	if cfg == nil {
		cfg = &Settings{}
	}

	viper.SetConfigType(configType)
	viper.SetEnvPrefix("SNEAK")
	viper.AutomaticEnv()
	viper.SetConfigName(configName)

	home, err := homedir.Dir()
	if err != nil {
		gui.ExitWithError(err)
	}

	cfgPath := fmt.Sprintf("%s/.sneak", home)

	// check for whether the directory and config file already exist
	err = fs.CreateDir(cfgPath)
	if err != nil {
		gui.ExitWithError(err)
	}

	viper.AddConfigPath(cfgPath)
	viper.Set("cfg_dir", cfgPath)

	if err := viper.ReadInConfig(); err != nil {
		gui.Warn(fmt.Sprintf("not seeing a config file where i'd expect it in %s - one sec...", cfgPath), nil)
	}
}

// SafeWriteConfig creates the config file if it doesn't exist yet
func SafeWriteConfig() error {
	// configFilePath := path.Join(viper.GetString())
	dir := viper.GetString("cfg_dir")
	filepath := path.Join(dir, configName+"."+configType)

	dirExists := fs.DirExists(dir)
	if !dirExists {
		err := fs.CreateDir(dir)
		if err != nil {
			return fmt.Errorf("could not create the config directory for sneak: %s", err)
		}
	}

	exists := fs.FileExists(filepath)
	if exists {
		return nil
	}

	gui.Info("popcorn", "creating your config file...", filepath)

	if _, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0600); err != nil {
		gui.Warn("could not create configuration file", filepath)
		return err
	}

	// set defaults
	htbUsername := gui.InputPromptWithResponse("what is your hack the box username?", "", true)
	viper.Set("htb_username", htbUsername)
	viper.Set("openvpn_filepath", fmt.Sprintf("%s/%s.ovpn", dir, viper.Get("htb_username")))
	preferredEditor := gui.GetUsersPreferredEditor("", true)
	viper.Set("default_editor", preferredEditor)
	viper.Set("data", dir)

	// write config file
	gui.Info("popcorn", "writing sneak defaults...", filepath)
	return viper.WriteConfigAs(filepath)
}

// Banner is the banner for sneak
var Banner = color.RedString(`

   ▄▄▄▄▄    ▄   ▄███▄   ██   █  █▀ 
  █     ▀▄   █  █▀   ▀  █ █  █▄█   
▄  ▀▀▀▀▄ ██   █ ██▄▄    █▄▄█ █▀▄   
 ▀▄▄▄▄▀  █ █  █ █▄   ▄▀ █  █ █  █  
         █  █ █ ▀███▀      █   █   
         █   ██           █   ▀    
                         ▀         														
`)
