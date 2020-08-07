package gui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

// printError outputs a standard error message to the user's console
func printError(msg interface{}) {
	color.HiRed(fmt.Sprintf("[ERROR %s]: %v", emoji.Sprint(":skull:"), msg))
}

// ExitWithError outputs a standard, formatted error and exits the program
func ExitWithError(msg interface{}) {
	printError(msg)
	os.Exit(1)
}
