package application

import (
	"fmt"
	"os"
	"os/user"

	comms "../commands"
)

/*RunCommands for Kanban*/
func (kanban *Kanban) RunCommands(args []string) {
	if len(args) == 1 {
		comms.Show(kanban.Storage)
	} else if len(args) == 2 {
		object := args[1]
		if "user" == object {
			showUser()
		} else {
			showHelp()
		}
	} else if len(args) == 3 {
		object := args[1]
		verb := args[2]
		comms.ExecuteOn(kanban.Storage, object, verb, "")
		comms.Show(kanban.Storage)
	} else if len(args) == 4 {
		object := args[1]
		verb := args[2]
		param := args[3]
		comms.ExecuteOn(kanban.Storage, object, verb, param)
		comms.Show(kanban.Storage)
	} else if len(args) > 4 {
		showHelp()
	}
}

func showHelp() {
	fmt.Println("Usage:", os.Args[0], "verb", "object")
	fmt.Println("Objects (verbs): lane (add, delete), card (add, delete, move), user")
}

func showUser() {
	currentuser, _ := user.Current()
	fmt.Println(currentuser.Name)
	fmt.Println(currentuser.HomeDir)
}
