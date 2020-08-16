package htb

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/internal/usecase/action"
	"github.com/kathleenfrench/sneak/pkg/file"
	"github.com/spf13/viper"
)

type executor struct {
	actionName string
	boxName    string
	runner     *entity.Runner
	fm         file.Manager
	box        entity.Box
	env        []string
	outputPath string
}

// RunPipeline handles running the scripts defined in every job in a pipeline
func (bg *BoxGUI) RunPipeline(p *entity.Pipeline) error {
	fmt.Println(fmt.Sprintf("%s %s", color.GreenString("[RUNNING PIPELINE]:"), p.Name))
	env := []string{
		fmt.Sprintf("BOX_IP=%s", bg.activeBox.IP),
		fmt.Sprintf("BOX_NAME=%s", bg.activeBox.Name),
		fmt.Sprintf("BOX_HOST=%s", bg.activeBox.Hostname),
	}

	actionUsecase := action.NewActionUsecase(bg.pipUsecase)
	for _, j := range p.Jobs {
		jobActions, err := actionUsecase.GetJobActions(j.Actions)
		if err != nil {
			return err
		}

		// handle runninger
		for _, a := range jobActions {
			exe := &executor{
				actionName: a.Name,
				runner:     a.Runner,
				fm:         file.NewManager(),
				box:        bg.activeBox,
				env:        env,
				outputPath: fmt.Sprintf("%s/notes/%s/main.md", viper.GetString("cfg_dir"), bg.activeBox.Name),
			}

			fmt.Println(
				fmt.Sprintf(
					"%s%s",
					color.YellowString("[ACTION]: "),
					a.Name),
			)

			err = exe.run(a.Name, a.Runner)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *executor) run(actionName string, r *entity.Runner) error {
	_, err := checkForNoteFile(e.box.Name)
	if err != nil {
		return err
	}

	switch r.ScriptPath {
	case "":
		return e.runCommand()
	default:
		return e.runScriptFile()
	}
}

func (e *executor) runScriptFile() error {
	path := e.runner.ScriptPath
	fileExists, err := e.fm.FileExists(path)
	if err != nil {
		return err
	}

	if !fileExists {
		return fmt.Errorf("the file %s does not exist", path)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	// verify executable
	if stat.Mode()&0100 != 0 {
		return fmt.Errorf("%s is not executable", path)
	}

	cmd := exec.Command("sh", path)
	cmd.Env = append(os.Environ(), e.env...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (e *executor) runCommand() (err error) {
	var stdout bytes.Buffer

	args := append([]string{"-c", e.runner.Command, "sh"})
	cmd := exec.Command("sh", args...)
	cmd.Env = append(os.Environ(), e.env...)
	cmd.Stdin = os.Stdin
	if e.runner.DontSaveLogs {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = &stdout
	}

	cmd.Stderr = os.Stderr
	gui.Spin.Start()
	err = cmd.Run()
	gui.Spin.Stop()
	if err != nil {
		return err
	}

	if !e.runner.DontSaveLogs {
		newDivider := fmt.Sprintf("\n\n## %s\n\n", e.actionName)
		err = e.fm.AppendToFile(e.outputPath, []byte(newDivider))
		if err != nil {
			return err
		}

		err = e.fm.AppendToFile(e.outputPath, stdout.Bytes())
		if err != nil {
			return err
		}

		fmt.Println(stdout.String())
	}

	gui.Info(":+1", fmt.Sprintf("logs written to %s", e.outputPath), nil)
	return nil
}
