package vpn

import (
	"net"

	"github.com/spf13/viper"
)

// OpenVPN represents the openvpn configs
type OpenVPN struct {
	Protocol     string
	Filepath     string
	LocalNetwork string
	IP           net.IP
	Port         uint16
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

// Connect establishes a connection with the vpn client
func (o *OpenVPN) Connect() error {
	return nil
}

// Same checks whether two connections are the same
func (o *OpenVPN) Same(conn OpenVPN) bool {
	return o.IP.Equal(conn.IP) && o.Port == conn.Port && o.Protocol == conn.Protocol
}
