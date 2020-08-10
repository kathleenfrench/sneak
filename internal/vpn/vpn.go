package vpn

import (
	"net"
	"os"
	"os/exec"

	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/common/shell"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// OpenVPN represents the openvpn configs
type OpenVPN struct {
	Protocol     string
	Filepath     string
	LocalNetwork string
	IP           net.IP
	Port         uint16
	home         string
	runner       shell.Runner
}

// NewOpenVPNClient returns a new wrapper for managing the openvpn client
func NewOpenVPNClient() (*OpenVPN, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	openvpn := &OpenVPN{
		Filepath: viper.GetString("openvpn_filepath"),
		runner:   shell.NewRunner(),
		home:     home,
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
	if !fs.FileExists("/tmp/ovpn_path") {
		err := fs.CreateFile("/tmp/ovpn_path")
		if err != nil {
			return err
		}

		uname, err := os.OpenFile("/tmp/ovpn_path", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer uname.Close()
		_, err = uname.WriteString(o.Filepath)
		if err != nil {
			return err
		}
	}

	cmd := &exec.Cmd{
		Path:   "/opt/vpn",
		Args:   []string{},
		Stdout: os.Stdout,
		Stderr: os.Stdin,
	}

	cmd.Start()
	cmd.Wait()

	return nil
}

// Same checks whether two connections are the same
func (o *OpenVPN) Same(conn OpenVPN) bool {
	return o.IP.Equal(conn.IP) && o.Port == conn.Port && o.Protocol == conn.Protocol
}
