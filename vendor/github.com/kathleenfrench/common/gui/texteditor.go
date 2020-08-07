package gui

import (
	"errors"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/exec"
)

// GetUsersPreferredEditor prompts the user for a text editor selection
func GetUsersPreferredEditor(label string, disableClearScreen bool) string {
	if label == "" {
		label = "what do you want to set as your default text editor?"
	}

	return SelectPromptWithResponse(label, DefaultEditors, "vim", disableClearScreen)
}

// DefaultEditors are some available text editors to set as defaults in applicable configs
var DefaultEditors = []string{
	"vim",
	"emacs",
	"vscode",
	"atom",
	"sublime",
	"phpstorm",
	"pstorm",
}

// EditorLaunchCommands are the commands used to open a file in a specified text editor and wait for the file to be saved to close
var EditorLaunchCommands = map[string]string{
	"vim":      "vim",
	"emacs":    "emacs",
	"vscode":   "code --wait",
	"atom":     "atom --wait",
	"sublime":  "subl -n -w",
	"phpstorm": "phpstorm --wait",
	"pstorm":   "pstorm --wait",
}

// GetPHPStormExecutable checks for which phpstorm executable is set in their PATH to determine the launch command to use
func GetPHPStormExecutable() (string, error) {
	color.HiBlue("checking for phpstorm executable...")
	phpstorm, err := exec.BashExec("phpstorm -v 2>/dev/null")
	if err != nil || phpstorm == "" {
		color.HiYellow("phpstorm executable not found")
		color.HiBlue("checking for pstorm executable...")
		pstorm, err := exec.BashExec("phpstorm -v 2>/dev/null")
		if err != nil || pstorm == "" {
			color.HiYellow("pstorm executable not found")
			return "", errors.New("it looks like you don't have the phpstorm command line installed")
		}

		return "pstorm", nil
	}

	return "phpstorm", nil
}
