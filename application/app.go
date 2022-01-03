package application

import (
	store "github.com/wlanboy/kanbantabs/v2/store"
)

/*Kanban containing workplace and storage*/
type Kanban struct {
	Storage  *store.Storage
	Filename string
}
