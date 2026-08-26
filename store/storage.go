package store

import (
	"encoding/json"
	"fmt"
	"os"

	model "github.com/wlanboy/kanbantabs/v2/model"
)

/*Storage containing everything*/
type Storage struct {
	Workplace *model.Workplace
	Filename  string
}

/*Load Workplace*/
func (storage *Storage) Load() {
	var workplace model.Workplace = model.Workplace{}

	file, err := os.ReadFile(storage.Filename)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println(err)
		}
		workplace.Name = "Kanban"
		workplace.NextID = 1
		storage.Workplace = &workplace
		storage.Save()
		return
	}

	json.Unmarshal(file, &workplace)
	storage.Workplace = &workplace
}

/*Save Workplace*/
func (storage *Storage) Save() {
	file, _ := json.MarshalIndent(storage.Workplace, "", " ")
	_ = os.WriteFile(storage.Filename, file, 0640)
}
