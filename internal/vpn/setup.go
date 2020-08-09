package vpn

import (
	"io/ioutil"

	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/common/gui"
)

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
