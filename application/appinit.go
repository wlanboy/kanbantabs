package application

import (
	"os/user"
	"path/filepath"

	"github.com/joho/godotenv"
	store "github.com/wlanboy/kanbantabs/v2/store"
)

/*Initialize configuration*/
func (kanban *Kanban) Initialize() {
	godotenv.Load()

	homedir := ""
	usr, err := user.Current()
	if err != nil {
		homedir = "/home/kanban"
	} else {
		homedir = usr.HomeDir
	}
	kanban.Filename = filepath.Join(homedir, ".kanbantabs")

	var storage store.Storage = store.Storage{}
	storage.Filename = kanban.Filename
	storage.Load()

	kanban.Storage = &storage
}
