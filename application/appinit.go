package application

import (
	store "../store"
	"github.com/joho/godotenv"
)

/*Initialize configuration*/
func (kanban *Kanban) Initialize() {
	godotenv.Load()

	kanban.Filename = "~/.kanbantabs"

	var storage store.Storage = store.Storage{}
	storage.Filename = kanban.Filename
	storage.Load()

	kanban.Storage = &storage
}
