package application

import (
	"os/user"
	"path/filepath"

	store "../store"
	"github.com/joho/godotenv"
)

/*Initialize configuration*/
func (kanban *Kanban) Initialize() {
	godotenv.Load()

	usr, _ := user.Current()
	kanban.Filename = filepath.Join(usr.HomeDir, ".kanbantabs")

	var storage store.Storage = store.Storage{}
	storage.Filename = kanban.Filename
	storage.Load()

	kanban.Storage = &storage
}
