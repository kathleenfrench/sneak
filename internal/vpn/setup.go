package vpn

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/kathleenfrench/common/exec"
	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/common/gui"
	"github.com/mitchellh/go-homedir"
)

// AlreadySetup is a helper to check whether or not a user's openvpn configs have already been setup
func (o *OpenVPN) AlreadySetup() bool {
	o.envCheck()

	if fs.FileExists(o.Filepath) && checkForPrivoxyConfig() {
		return true
	}

	return false
}

func (o *OpenVPN) envCheck() {
	openVPNConfigPath := os.Getenv("OPENVPN_CONF")
	if len(openVPNConfigPath) == 0 {
		os.Setenv("OPENVPN_CONF", o.Filepath)
	}
}

// Setup creates the openvpn file at the expected location for sneak and prompts the user for its contents, then writes them to the file
func (o *OpenVPN) Setup(defaultEditor string) error {
	o.envCheck()

	err := o.createConfigFile()
	if err != nil {
		return err
	}

	vpnCfgs := gui.TextEditorInputAndSave("copy in your openvpn file from HTB and save it", "", defaultEditor)
	err = ioutil.WriteFile(o.Filepath, []byte(vpnCfgs), 0644)
	if err != nil {
		gui.ExitWithError(err)
	}

	// create privoxy config from default
	err = createFunctionalPrivoxyConfig()
	if err != nil {
		return err
	}

	return nil
}

func createFunctionalPrivoxyConfig() error {
	whoami, err := exec.BashExec("whoami")
	if err != nil {
		return err
	}

	// need to include the confdir and logdir lines in the privoxy config file for it to work
	appendLines := fmt.Sprintf("\n%s\n%s\n", fmt.Sprintf("confdir /home/%s", whoami), fmt.Sprintf("logdir /home/%s", whoami))

	home, err := homedir.Dir()
	if err != nil {
		gui.ExitWithError(err)
	}

	defaultConfigs, err := os.Open(fmt.Sprintf("%s/config.default", home))
	if err != nil {
		return fmt.Errorf("could not locate the default privoxy configs - %s", err)
	}

	defer defaultConfigs.Close()

	privoxyConfig, err := os.Create(fmt.Sprintf("%s/config", home))
	if err != nil {
		return err
	}

	defer privoxyConfig.Close()

	_, err = io.Copy(privoxyConfig, defaultConfigs)
	if err != nil {
		return fmt.Errorf("could not write default privoxy file to new config file - %s", err)
	}

	privoxyConfigPreModify, err := os.OpenFile(fmt.Sprintf("%s/config", home), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer privoxyConfigPreModify.Close()
	_, err = privoxyConfigPreModify.WriteString(appendLines)
	if err != nil {
		return err
	}

	return nil
}

func checkForPrivoxyConfig() bool {
	home, err := homedir.Dir()
	if err != nil {
		gui.ExitWithError(err)
	}

	privoxyConfig := fmt.Sprintf("%s/config", home)
	if fs.FileExists(privoxyConfig) {
		return true
	}

	return false
}
