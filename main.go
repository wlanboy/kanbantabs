package main

import (
	"os"

	app "./application"
)

func main() {
	kanban := app.Kanban{}
	kanban.Initialize()
	kanban.RunCommands(os.Args)
}
