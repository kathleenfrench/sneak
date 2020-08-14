package htb

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/pkg/utils"
)

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

	utils.Spacer()
	gui.SideBySideTable(data, "Red")
	utils.Spacer()
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

	utils.Spacer()
	gui.SideBySideTable(data, "Red")
	utils.Spacer()
}
