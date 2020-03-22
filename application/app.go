package application

import (
	store "../store"
)

/*Kanban containing workplace and storage*/
type Kanban struct {
	Storage  *store.Storage
	Filename string
}
