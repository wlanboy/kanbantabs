package store

import (
	"strconv"

	model "github.com/wlanboy/kanbantabs/v2/model"
)

/*AddBoard to Workplace*/
func (storage *Storage) AddBoard(board model.Board) {
	storage.Workplace.Lanes = append(storage.Workplace.Lanes, board)
	storage.Save()
}

/*DeleteBoard to Workplace*/
func (storage *Storage) DeleteBoard(boardnumber string) {
	number, err := strconv.ParseInt(boardnumber, 10, 32)
	if err == nil {
		index := int(number) - 1

		if index >= 0 && index < len(storage.Workplace.Lanes) {
			storage.Workplace.Lanes[index] = storage.Workplace.Lanes[len(storage.Workplace.Lanes)-1]
			storage.Workplace.Lanes = storage.Workplace.Lanes[:len(storage.Workplace.Lanes)-1]
			storage.Save()
		}
	}
}
