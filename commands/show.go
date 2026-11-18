package commands

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	model "github.com/wlanboy/kanbantabs/v2/model"
	store "github.com/wlanboy/kanbantabs/v2/store"
)

const columnWidth = 25

/*Show prints all lanes side by side, one card per row*/
func Show(storage *store.Storage, out io.Writer) {
	lanes := storage.Workplace.Lanes
	if len(lanes) == 0 {
		fmt.Fprintln(out, "No lanes yet. Add one with: kanbantabs lane add <name>")
		return
	}

	rows := 0
	var header strings.Builder
	for number, lane := range lanes {
		fmt.Fprintf(&header, "%d: %s\t", number+1, lane.Name)
		rows = max(rows, len(lane.Items))
	}

	writer := tabwriter.NewWriter(out, columnWidth, columnWidth, 0, ' ', tabwriter.Debug)
	separator := strings.Repeat("-", len(lanes)*(columnWidth+1))

	fmt.Fprintln(writer, separator)
	fmt.Fprintln(writer, header.String())
	fmt.Fprintln(writer, separator)
	for row := range rows {
		var line strings.Builder
		for _, lane := range lanes {
			if row < len(lane.Items) {
				line.WriteString(formatCard(lane.Items[row]))
			}
			line.WriteString("\t")
		}
		fmt.Fprintln(writer, line.String())
	}
	writer.Flush()
}

func formatCard(card model.BoardItem) string {
	return fmt.Sprintf("[%d%s] %s", card.ID, card.Severity.Marker(), card.Name)
}
