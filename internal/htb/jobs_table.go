package htb

import (
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/pkg/utils"
)

func printJobTable(job *entity.Job) {
	data := []table.Row{
		{"name", job.Name},
		{"description", job.Description},
		{"actions", strings.Join(job.Actions, ",")},
		{"disabled", job.Disabled},
	}

	utils.Spacer()
	gui.SideBySideTable(data, "Red")
	utils.Spacer()
}
