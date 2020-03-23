package store

import (
	"strconv"

	model "../model"
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
		number := int(number) + 1

		if number > 0 && number < len(storage.Workplace.Lanes) {
			storage.Workplace.Lanes[number] = storage.Workplace.Lanes[len(storage.Workplace.Lanes)-1]
			storage.Workplace.Lanes = storage.Workplace.Lanes[:len(storage.Workplace.Lanes)-1]
			storage.Save()
		}
	}
}
