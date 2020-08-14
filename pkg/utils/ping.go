package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// SudoPing becasuse the alpine linux docker imnage does not play nice with normal ping
func SudoPing(ip string) error {
	cmd := exec.Command("sudo", "ping", "-c", fmt.Sprintf("%d", 1), ip)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
