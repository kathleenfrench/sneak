package htb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/kathleenfrench/sneak/internal/helpers"
	"github.com/spf13/viper"

	humanize "github.com/dustin/go-humanize"
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
func PromptUserForBoxData() (Box, error) {
	box := Box{
		Name:       gui.InputPromptWithResponse("what is the name of the box?", "", true),
		IP:         gui.InputPromptWithResponse("what is its IP?", "", true),
		Completed:  false,
		Active:     false,
		Notes:      "",
		OS:         gui.SelectPromptWithResponse("what is the OS?", osOptions, nil, true),
		Difficulty: gui.SelectPromptWithResponse("what is its difficulty?", difficulties, nil, true),
		Up:         false,
		Flags: Flags{
			Root: "",
			User: "",
		},
		Created:     time.Now(),
		LastUpdated: time.Now(),
	}

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

// CompletionColorizer returns an icon with the completion status of a box
func CompletionColorizer(completed bool) string {
	if completed {
		return color.HiGreenString("pwnd")
	}

	return color.HiYellowString("incomplete")
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

func constructBoxListing(box Box) string {
	head := fmt.Sprintf(
		"%s - [%s][%s][%s]",
		box.Name,
		color.HiBlueString(box.OS),
		DifficultyColorizer(box.Difficulty),
		CompletionColorizer(box.Completed),
	)

	return head
}

func makeGuiBoxMappings(boxes []Box) (keys []string, mapping map[string]Box) {
	mapping = make(map[string]Box)

	for _, b := range boxes {
		name := constructBoxListing(b)
		keys = append(keys, name)
		mapping[name] = b
	}

	return keys, mapping
}

// SelectBoxFromDropdown lists a collection of boxes to choose from in a terminal dropdown
func SelectBoxFromDropdown(boxes []Box) Box {
	boxNames, boxMap := makeGuiBoxMappings(boxes)
	selection := gui.SelectPromptWithResponse("select a box", boxNames, nil, false)
	selected := boxMap[selection]
	return selected
}

// PrintBoxDataTable poutputs box data in a readable table in the terminal window
func PrintBoxDataTable(box Box) {
	data := []table.Row{
		{"name", box.Name},
		{"IP", box.IP},
		{"description", box.Description},
		{"hostname", box.Hostname},
		{"os", box.Hostname},
		{"difficulty", box.Difficulty},
		{"added", humanize.Time(box.Created)},
		{"last updated", humanize.Time(box.LastUpdated)},
	}

	helpers.Spacer()
	gui.SideBySideTable(data, "HiRed")
	helpers.Spacer()
}

const (
	setToActive     = "set to active"
	checkConnection = "check connection"
	openNotes       = "open notes"
	flags           = "flags"
	returnToBoxes   = "return to other boxes"
	quit            = "quit"
)

var boxActions = []string{
	setToActive,
	checkConnection,
	openNotes,
	flags,
	returnToBoxes,
	quit,
}

// SelectBoxActionsDropdown lists available actions with a single box or the ability to return to the 'main menu' of boxes
func SelectBoxActionsDropdown(box Box, boxes []Box) error {
	PrintBoxDataTable(box)
	selection := gui.SelectPromptWithResponse("select from the dropdown", boxActions, nil, true)

	switch selection {
	case setToActive:
		color.Red("TODO")
	case checkConnection:
		color.Red("TODO")
	case openNotes:
		note, err := checkForNoteFile(box.Name)
		if err != nil {
			return err
		}

		updatedNote := gui.TextEditorInputAndSave(fmt.Sprintf("update your notes on %s in markdown", box.Name), note, viper.GetString("default_editor"))

		err = saveNoteFile(box.Name, updatedNote)
		if err != nil {
			return err
		}

		return SelectBoxActionsDropdown(box, boxes)
	case flags:
		color.Red("TODO")
	case returnToBoxes:
		return SelectBoxActionsDropdown(SelectBoxFromDropdown(boxes), boxes)
	case quit:
		os.Exit(0)
	}

	return nil
}

func saveNoteFile(boxName string, note string) error {
	notesPath := config.GetNotesDirectory()
	notesFilePath := fmt.Sprintf("%s/%s.md", notesPath, boxName)

	err := ioutil.WriteFile(notesFilePath, []byte(note), 0644)
	if err != nil {
		return err
	}

	return nil
}

func checkForNoteFile(boxName string) (string, error) {
	notesPath := config.GetNotesDirectory()
	notesFilePath := fmt.Sprintf("%s/%s.md", notesPath, boxName)

	// if the notes file already exists, read the text from the file and return it as a string to set a s adefault
	if fs.FileExists(notesFilePath) {
		note, err := ioutil.ReadFile(notesFilePath)
		if err != nil {
			return "", err
		}

		return string(note), nil
	}

	err := fs.CreateFile(notesFilePath)
	if err != nil {
		return "", err
	}

	return "", nil
}
