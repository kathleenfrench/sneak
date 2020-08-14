package htb

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/helpers"
	"github.com/kathleenfrench/sneak/internal/usecase/boxusecase"
	"github.com/spf13/viper"

	humanize "github.com/dustin/go-humanize"
)

// BoxGUI is an interface for methods managing the box GUI
type BoxGUI interface {
	SelectBoxFromDropdown(boxes []entity.Box) entity.Box
	SelectBoxActionsDropdown(box entity.Box, boxes []entity.Box) error
	PromptUserForBoxData() (entity.Box, error)
}

type boxGUI struct {
	singleBoxTableShown bool
	activeBox           string
	usecase             boxusecase.Usecase
}

// NewBoxGUI instantiates a new box gui interface
func NewBoxGUI(use boxusecase.Usecase) BoxGUI {
	return &boxGUI{
		singleBoxTableShown: false,
		activeBox:           "",
		usecase:             use,
	}
}

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
func (bg *boxGUI) PromptUserForBoxData() (entity.Box, error) {
	box := entity.Box{
		Name:        gui.InputPromptWithResponse("what is the name of the box?", "", true),
		IP:          gui.InputPromptWithResponse("what is its IP?", "", true),
		Description: gui.InputPromptWithResponse("provide a short description of the box if you want", "", true),
		Completed:   false,
		Active:      false,
		Notes:       "",
		OS:          gui.SelectPromptWithResponse("what is the OS?", osOptions, nil, true),
		Difficulty:  gui.SelectPromptWithResponse("what is its difficulty?", difficulties, nil, true),
		Up:          false,
		Flags: entity.Flags{
			Root: "",
			User: "",
		},
		Created:     time.Now(),
		LastUpdated: time.Now(),
	}

	if err := box.Validate(); err != nil {
		return box, err
	}

	box.Hostname = fmt.Sprintf("%s.htb", box.Name)
	return box, nil
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

func constructBoxListing(box entity.Box) string {
	head := fmt.Sprintf(
		"%s - [%s][%s][%s]",
		box.Name,
		color.HiBlueString(box.OS),
		DifficultyColorizer(box.Difficulty),
		CompletionColorizer(box.Completed),
	)

	return head
}

func makeGuiBoxMappings(boxes []entity.Box) (keys []string, mapping map[string]entity.Box) {
	mapping = make(map[string]entity.Box)

	for _, b := range boxes {
		name := constructBoxListing(b)
		keys = append(keys, name)
		mapping[name] = b
	}

	return keys, mapping
}

// SelectBoxFromDropdown lists a collection of boxes to choose from in a terminal dropdown
func (bg *boxGUI) SelectBoxFromDropdown(boxes []entity.Box) entity.Box {
	boxNames, boxMap := makeGuiBoxMappings(boxes)
	selection := gui.SelectPromptWithResponse("select a box", boxNames, nil, false)
	selected := boxMap[selection]
	return selected
}

// PrintBoxDataTable poutputs box data in a readable table in the terminal window
func PrintBoxDataTable(box entity.Box) {
	data := []table.Row{
		{"name", box.Name},
		{"IP", box.IP},
		{"description", box.Description},
		{"hostname", box.Hostname},
		{"os", box.OS},
		{"difficulty", box.Difficulty},
		{"active", box.Active},
		{"completed", box.Completed},
		{"added", humanize.Time(box.Created)},
		{"last updated", humanize.Time(box.LastUpdated)},
	}

	helpers.Spacer()
	gui.SideBySideTable(data, "Red")
	helpers.Spacer()
}

func printFlagTable(flags entity.Flags) {
	userFlag := flags.User
	rootFlag := flags.Root

	if userFlag == "" {
		userFlag = "NOT SET"
	}

	if rootFlag == "" {
		rootFlag = "NOT SET"
	}

	data := []table.Row{
		{"user", userFlag},
		{"root", rootFlag},
	}

	helpers.Spacer()
	gui.SideBySideTable(data, "Red")
	helpers.Spacer()
}

const (
	toggleActiveStatus = "toggle active status"
	checkConnection    = "check connection"
	openNotes          = "open notes editor"
	quickViewNotes     = "quickview notes"
	editDescription    = "edit description"
	flags              = "flags"
	returnToBoxes      = "return to other boxes"
	quit               = "quit"
	seeTable           = "show info table"
)

var boxActions = []string{
	toggleActiveStatus,
	checkConnection,
	openNotes,
	quickViewNotes,
	editDescription,
	flags,
	returnToBoxes,
	quit,
}

// SelectBoxActionsDropdown lists available actions with a single box or the ability to return to the 'main menu' of boxes
func (bg *boxGUI) SelectBoxActionsDropdown(box entity.Box, boxes []entity.Box) error {
	if box.Active {
		bg.activeBox = box.Name
	}

	if !bg.singleBoxTableShown {
		PrintBoxDataTable(box)
		boxActions = append([]string{seeTable}, boxActions...)
	}

	bg.singleBoxTableShown = true
	selection := gui.SelectPromptWithResponse("select from the dropdown", boxActions, nil, true)

	switch selection {
	case seeTable:
		PrintBoxDataTable(box)
		return bg.SelectBoxActionsDropdown(box, boxes)
	case toggleActiveStatus:
		switch box.Active {
		case true:
			color.HiGreen("%s is currently set to active", box.Name)
			setInactive := gui.ConfirmPrompt(fmt.Sprintf("set %s as inactive?", box.Name), "", true, true)
			switch setInactive {
			case true:
				box.Active = false
				bg.activeBox = ""
			default:
				return bg.SelectBoxActionsDropdown(box, boxes)
			}
		default:
			color.Red("%s is not currently set to active", box.Name)
			setActive := gui.ConfirmPrompt(fmt.Sprintf("set %s as active?", box.Name), "", true, true)
			switch setActive {
			case true:
				box.Active = true
				bg.activeBox = box.Name
			default:
				return bg.SelectBoxActionsDropdown(box, boxes)
			}
		}

		// write the change to the db
		err := bg.usecase.Save(box)
		if err != nil {
			return err
		}

		color.Green("successfully changed the active status of %s!", box.Name)

		// after making that change, re-fetch all of our boxes for up to date info
		boxes, err = bg.usecase.GetAll()
		if err != nil {
			return err
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
	case quickViewNotes:
		note, err := checkForNoteFile(box.Name)
		if err != nil {
			return err
		}

		if len(note) == 0 {
			color.Yellow("you have not started a note for %s yet!", box.Name)
		} else {
			fmt.Println(helpers.RenderMarkdown(note))
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
	case checkConnection:
		err := helpers.SudoPing(box.IP)
		if err != nil {
			gui.Warn("uh oh, that box couldn't be reached! verify that the machine is active and your VPN connection is still intact", box.IP)
		} else {
			gui.Info("+1", "reachable!", box.IP)
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
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

		return bg.SelectBoxActionsDropdown(box, boxes)
	case editDescription:
		box.Description = gui.InputPromptWithResponse("provide a new description", "", true)
		// write the change to the db
		err := bg.usecase.Save(box)
		if err != nil {
			return err
		}

		// after making that change, re-fetch all of our boxes for up to date info
		boxes, err = bg.usecase.GetAll()
		if err != nil {
			return err
		}

		return bg.SelectBoxActionsDropdown(box, boxes)
	case flags:
		printFlagTable(box.Flags)
		addOrUpdate := gui.ConfirmPrompt("do you want to update any flag values?", "", false, true)
		switch addOrUpdate {
		case true:
			whichFlags := gui.SelectPromptWithResponse("which flag?", []string{"root", "user"}, nil, true)
			switch whichFlags {
			case "root":
				box.Flags.Root = gui.InputPromptWithResponse("enter a new root flag", box.Flags.Root, true)
			case "user":
				box.Flags.User = gui.InputPromptWithResponse("enter a new user flag", box.Flags.User, true)
			}

			// write the change to the db
			err := bg.usecase.Save(box)
			if err != nil {
				return err
			}

			// after making that change, re-fetch all of our boxes for up to date info
			boxes, err = bg.usecase.GetAll()
			if err != nil {
				return err
			}

			fallthrough
		default:
			return bg.SelectBoxActionsDropdown(box, boxes)
		}
	case returnToBoxes:
		return bg.SelectBoxActionsDropdown(bg.SelectBoxFromDropdown(boxes), boxes)
	case quit:
		os.Exit(0)
	}

	return nil
}

func saveNoteFile(boxName string, note string) error {
	notesPath := fmt.Sprintf("%s/%s", config.GetNotesDirectory(), boxName)

	// create note directory for that box if it doesn't exist
	err := fs.CreateDir(notesPath)
	if err != nil {
		return fmt.Errorf("there was an error creating the notes directory for %s - %w", boxName, err)
	}

	notesFilePath := fmt.Sprintf("%s/main.md", notesPath)
	err = ioutil.WriteFile(notesFilePath, []byte(note), 0644)
	if err != nil {
		return err
	}

	return nil
}

func checkForNoteFile(boxName string) (string, error) {
	notesPath := fmt.Sprintf("%s/%s", config.GetNotesDirectory(), boxName)
	notesFilePath := fmt.Sprintf("%s/main.md", notesPath)

	// if the notes file already exists, read the text from the file and return it as a string to set a s adefault
	if fs.FileExists(notesFilePath) {
		note, err := ioutil.ReadFile(notesFilePath)
		if err != nil {
			return "", err
		}

		return string(note), nil
	}

	// create the directory if it doesn't exist yet
	err := fs.CreateDir(notesPath)
	if err != nil {
		return "", err
	}

	err = fs.CreateFile(notesFilePath)
	if err != nil {
		return "", err
	}

	return "", nil
}
