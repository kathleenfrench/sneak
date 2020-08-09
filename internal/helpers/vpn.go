package helpers

import (
	"io/ioutil"

	"github.com/kathleenfrench/common/exec"
	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/common/gui"
	"github.com/spf13/viper"
)

// OpenVPN represents the openvpn configs
type OpenVPN struct {
	Filepath     string
	LocalNetwork string
}

// NewOpenVPNClient returns a new wrapper for managing the openvpn client
func NewOpenVPNClient() (*OpenVPN, error) {
	openvpn := &OpenVPN{
		Filepath: viper.GetString("openvpn_filepath"),
	}

	ln, err := localNetwork()
	if err != nil {
		return nil, err
	}

	openvpn.LocalNetwork = ln

	return openvpn, nil
}

func localNetwork() (string, error) {
	ip, err := exec.BashExec(`ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}'`)
	if err != nil {
		return ip, err
	}

	return ip, nil
}

// Connect establishes a connection with the vpn client
func (o *OpenVPN) Connect() error {
	return nil
}

// AlreadySetup is a helper to check whether or not a user's openvpn configs have already been setup
func (o *OpenVPN) AlreadySetup() bool {
	if fs.FileExists(o.Filepath) {
		return true
	}

	return false
}

// Setup creates the openvpn file at the expected location for sneak and prompts the user for its contents, then writes them to the file
func (o *OpenVPN) Setup(defaultEditor string) error {
	err := o.createConfigFile()
	if err != nil {
		return err
	}

	vpnCfgs := gui.TextEditorInputAndSave("copy in your openvpn file from HTB and save it", "", defaultEditor)
	err = ioutil.WriteFile(o.Filepath, []byte(vpnCfgs), 0644)
	if err != nil {
		gui.ExitWithError(err)
	}

	return nil
}

// CreateConfigFile writes the openvpn config file to the expected path for sneak
func (o *OpenVPN) createConfigFile() error {
	err := fs.CreateFile(o.Filepath)
	if err != nil {
		return nil
	}

	return nil
}
