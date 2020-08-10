package vpn

import (
	"github.com/kathleenfrench/common/exec"
	"github.com/kathleenfrench/common/fs"
)

// CreateConfigFile writes the openvpn config file to the expected path for sneak
func (o *OpenVPN) createConfigFile() error {
	err := fs.CreateFile(o.Filepath)
	if err != nil {
		return nil
	}

	return nil
}

func localNetwork() (string, error) {
	ip, err := exec.BashExec(`ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}'`)
	if err != nil {
		return ip, err
	}

	return ip, nil
}
