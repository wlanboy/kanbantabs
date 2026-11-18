package main

import (
	"os"

	app "github.com/wlanboy/kanbantabs/v2/application"
)

func main() {
	kanban := app.Kanban{}
	if err := kanban.Initialize(); err != nil {
		app.PrintError(err)
		os.Exit(1)
	}
	if err := kanban.Run(os.Args[1:]); err != nil {
		app.PrintError(err)
		os.Exit(1)
	}
}
