package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	model "../model"
	store "../store"
)

/*Show for Kanban*/
func Show(storage *store.Storage) {
	const padding = 0
	const length = 25
	tabwriter := tabwriter.NewWriter(os.Stdout, length, length, padding, ' ', tabwriter.TabIndent|tabwriter.Debug)

	lanes := storage.Workplace.Lanes
	numberOfLanes := len(lanes)
	currentLine := 1
	moreCards := true
	currentLineString := ""
	header := ""

	lineseparator := strings.Repeat("-", numberOfLanes*(length+1))
	currentLines := []string{}

	for moreCards == true {
		addedCards := 0
		currentLineString = ""

		for _, board := range lanes {
			if currentLine == 1 {
				header = header + writeHeader(&board)
			}

			if len(board.Items) >= currentLine {
				item := board.Items[currentLine-1]
				currentLineString = currentLineString + writeLane(&item)
				addedCards++
			} else {
				currentLineString = writeEmpty(currentLineString)
			}

		}

		if addedCards == 0 {
			moreCards = false
		} else {
			currentLines = append(currentLines, currentLineString)
			currentLine++
		}
	}

	if len(header) > 1 {
		fmt.Fprintln(tabwriter, lineseparator)
		fmt.Fprintln(tabwriter, header)
		fmt.Fprintln(tabwriter, lineseparator)

		for _, line := range currentLines {
			fmt.Fprintln(tabwriter, line)
		}
		tabwriter.Flush()
	} else {
		fmt.Println("Empty boards")
	}
}

func writeEmpty(line string) string {
	return fmt.Sprintf("%s\t", line)
}

func writeHeader(board *model.Board) string {
	return fmt.Sprintf("%s\t", board.Name)
}

func writeLane(item *model.BoardItem) string {
	return fmt.Sprintf("[%d] %s\t", item.ID, item.Name)
}
