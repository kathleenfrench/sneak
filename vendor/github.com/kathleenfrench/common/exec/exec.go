package exec

import (
	"os/exec"
	"strings"
)

// BashExec executes a command appended to a 'bash -c' command
func BashExec(cmd string) (string, error) {
	r, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	trimmed := strings.TrimSpace(string(r))
	return trimmed, nil
}
