package htb

import (
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/pkg/utils"
)

// printRunnerTable outputs a table with information about a given runner
func printRunnerTable(runner *entity.Runner) {
	var (
		scriptPath string
		outputPath string
		commandVal string
	)

	switch runner.Command {
	case "":
		commandVal = "n/a"
	default:
		commandVal = strings.TrimPrefix(runner.Command, "|")
	}

	switch runner.OutputPath {
	case "":
		outputPath = "uses box default"
	default:
		outputPath = runner.OutputPath
	}

	switch runner.ScriptPath {
	case "":
		scriptPath = "n/a"
	default:
		scriptPath = runner.ScriptPath
	}

	data := []table.Row{
		{"command", commandVal},
		{"script path", scriptPath},
		{"output path", outputPath},
	}

	utils.Spacer()
	gui.SideBySideTable(data, "Red")
	utils.Spacer()
}
