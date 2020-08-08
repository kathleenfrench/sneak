package htb

import (
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kyokomi/emoji"
)

var osOptions = []string{
	"linux",
	"windows",
	"freeBSD",
	"openBSD",
	"other",
}

var difficulties = []string{
	"easy",
	"medium",
	"hard",
	"insane",
}

// PromptUserForBoxData prompts the user for values about the htb machine they want to add
func PromptUserForBoxData() (*Box, error) {
	box := &Box{
		Name:      gui.InputPromptWithResponse("what is the name of the box?", "", true),
		IP:        gui.InputPromptWithResponse("what is its IP?", "", true),
		Completed: false,
		Active:    false,
		Notes:     "",
		Up:        false,
		Flags: Flags{
			Root: "",
			User: "",
		},
		Created:     time.Now(),
		LastUpdated: time.Now(),
	}

	os := gui.SelectPromptWithResponse("what is the OS?", osOptions, "linux", true)
	difficulty := gui.SelectPromptWithResponse("what is its difficulty?", difficulties, "easy", true)

	box.OS = os
	box.Difficulty = difficulty

	color.Red("difficulty: %s", box.Difficulty)
	color.Red("OS: %s", box.OS)

	if err := box.validate(); err != nil {
		return box, err
	}

	box.Hostname = fmt.Sprintf("%s.htb", box.Name)
	return box, nil
}

func (b Box) validate() error {
	if b.Name == "" {
		return errors.New("setting a name for the box is required")
	}

	if !govalidator.IsIP(b.IP) {
		gui.Warn("invalid IP address", b.IP)
		b.IP = gui.InputPromptWithResponse("what is its IP?", "", true)
		if !govalidator.IsIP(b.IP) {
			return errors.New(b.IP + " is not a valid IP address")
		}
	}

	return nil
}

// CompletionStatusIcon returns an icon with the completion status of a box
func CompletionStatusIcon(completed bool) string {
	if completed {
		return fmt.Sprintf("completed: %s", emoji.Sprint(":white_check_mark:"))
	}

	return fmt.Sprintf("completed: %s", emoji.Sprint(":x:"))
}

// DifficultyColorizer colorizes based on difficulty
func DifficultyColorizer(diff string) string {
	switch diff {
	case "easy":
		return color.GreenString("easy")
	case "medium":
		return color.YellowString("medium")
	case "hard":
		return color.HiRedString("hard")
	case "insane":
		return color.RedString("insane")
	}

	return ""
}
