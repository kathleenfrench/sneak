package helpers

import "github.com/kathleenfrench/common/exec"

import "github.com/spf13/viper"

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
