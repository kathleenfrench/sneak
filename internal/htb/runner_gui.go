package htb

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/action"
	"github.com/spf13/viper"
)

// RunnerGUI handles the terminal GUI for defining runners, e.g. scripts and tools to run for various pentesting operations, in sneak
type RunnerGUI struct {
	usecase action.Usecase
	ActionsGUIHandlers
}

// NewRunnerGUI instantiates a new GUI for defining runners
func NewRunnerGUI(au action.Usecase) *RunnerGUI {
	return &RunnerGUI{
		usecase: au,
	}
}

// existing runner opts
const (
	editCommand    = "change command(s) to run"
	editScriptPath = "change script path"
	editOutputPath = "change output path"
	returnToAction = "return to action"
)

// new runner opts
const (
	commandConfig     = "do you want to add a simple command or script directly, or point to a script saved at a particular filepath?"
	addCommand        = "add an executable command or script"
	addScriptPath     = "provide a filepath for the runner script"
	outputConfig      = "do you want to set a custom output path for the script logs or keep the deefault?"
	addOutputPath     = "set a custom filepath"
	keepOutputDefault = "save action logs to the default filepath in the context of a box"
)

var singleRunnerAddOpts = []string{
	commandConfig,
	outputConfig,
}

var commandChoiceOpts = []string{
	addCommand,
	addScriptPath,
}

var outputChoiceOpts = []string{
	keepOutputDefault,
	addOutputPath,
}

var singleRunnerOpts = []string{
	editCommand,
	editScriptPath,
	editOutputPath,
	returnToAction,
	quit,
}

var newRunnerOutputPathMsg = color.YellowString(`
[OUTPUT PATHS]: 

by default, sneak saves the output/logs for a runner in the notes/[box_name]
directory. if you would prefer to define a custom path for logs on this
action, you will need to provide one.
`)

var newRunnerDirections = color.RedString(`
[IMPORTANT]: when adding a new command or script to sneak, 
you must denote expected values in the following way:

[IP ADDRESSES]: indicate the box IP with $BOX_IP
[BOX NAMES]: indicate the box name with $BOX_NAME
[BOX HOSTNAMES]: indicate the box hostname with $BOX_HOSTNAME

be sure to define what type of script is running, i.e. begin your
file with #!/bin/bash or, if another language is to be used (and is installed)
to your environment)

sneak will handle parsing those values using what you have set
for that box.
`)

var defaultScriptText = fmt.Sprintf("#!/bin/bash\n\n")

// AddNewRunner manages collectin information on the addition of a new runner
func (rg *RunnerGUI) AddNewRunner() (*entity.Runner, error) {
	newRunner := &entity.Runner{}
	commandConfigChoice := gui.SelectPromptWithResponse(commandConfig, commandChoiceOpts, nil, false)
	fmt.Println(newRunnerDirections)
	switch commandConfigChoice {
	case addCommand:
		newRunner.Command = gui.TextEditorInputAndSave("write your script, then save and exit", defaultScriptText, viper.GetString("default_editor"))
		newRunner.Command = fmt.Sprintf("|\n%s", newRunner.Command)
	case addScriptPath:
		newRunner.ScriptPath = gui.InputPromptWithResponse("provide a path to the script", "", true)
	}

	outputConfigChoice := gui.SelectPromptWithResponse(outputConfig, outputChoiceOpts, keepOutputDefault, true)
	fmt.Println(newRunnerOutputPathMsg)
	switch outputConfigChoice {
	case addOutputPath:
		newRunner.OutputPath = gui.InputPromptWithResponse("provide a custom output path", "", true)
	case keepOutputDefault:
		break
	}

	return newRunner, nil
}

// HandleRunnerDropdown handles the actions available through the GUI dropdown for working with runners
func (rg *RunnerGUI) HandleRunnerDropdown(action *entity.Action) error {
	runnerOpts := gui.SelectPromptWithResponse("select from dropdown", singleRunnerOpts, nil, true)

	switch runnerOpts {
	case editCommand:
		action.Runner.Command = gui.TextEditorInputAndSave("write your script, then save and exit", action.Runner.Command, viper.GetString("default_editor"))
		action.Runner.Command = fmt.Sprintf("|\n%s", action.Runner.Command)
	case editScriptPath:
		action.Runner.ScriptPath = gui.InputPromptWithResponse("provide a path to the script", action.Runner.ScriptPath, true)
	case editOutputPath:
		action.Runner.OutputPath = gui.InputPromptWithResponse("provide a custom output path", action.Runner.OutputPath, true)
	case returnToAction:
		return rg.SelectIndividualActionsActionsDropdown(action)
	case quit:
		os.Exit(0)
	}

	err := rg.usecase.SaveAction(action)
	if err != nil {
		return err
	}

	return nil
}
