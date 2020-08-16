package htb

import (
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/kathleenfrench/sneak/internal/entity"
	"github.com/kathleenfrench/sneak/pkg/utils"
)

func printJobsTable(jobs map[string]*entity.Job) {
	var rows []table.Row

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)
	t.Style().Color.Header = text.Colors{text.BgBlack, text.FgWhite}
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Colors:   text.Colors{text.FgBlack},
			WidthMax: 60,
		},
	})
	t.SetOutputMirror(os.Stdout)
	header := table.Row{"NAME", "DESCRIPTION", "ACTIONS", "DIABLED"}
	t.AppendHeader(header)

	for j, data := range jobs {
		rows = append(rows, table.Row{j, data.Description, strings.Join(data.Actions, ","), data.Disabled})
	}

	t.AppendRows(rows)

	utils.Spacer()
	t.Render()
	utils.Spacer()
}
