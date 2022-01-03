package main

import (
	"os"

	app "github.com/wlanboy/kanbantabs/v2/application"
)

func main() {
	kanban := app.Kanban{}
	kanban.Initialize()
	kanban.RunCommands(os.Args)
}
