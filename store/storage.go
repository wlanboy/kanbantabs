package store

import (
	"encoding/json"
	"io/ioutil"
	"os"

	model "../model"
)

/*Storage containing everything*/
type Storage struct {
	Workplace *model.Workplace
	Filename  string
}

/*Load Workplace*/
func (storage *Storage) Load() {
	var workplace model.Workplace = model.Workplace{}

	info, _ := os.Stat(storage.Filename)
	if info != nil {
		file, _ := ioutil.ReadFile(storage.Filename)
		json.Unmarshal(file, &workplace)
		storage.Workplace = &workplace
	} else {
		workplace.Name = "Kanban"
		workplace.NextID = 1
		storage.Workplace = &workplace
		storage.Save()
	}
}

/*Save Workplace*/
func (storage *Storage) Save() {
	file, _ := json.MarshalIndent(storage.Workplace, "", " ")
	_ = ioutil.WriteFile(storage.Filename, file, 0640)
}
