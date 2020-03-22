package store

import (
	"fmt"
	"strconv"

	model "../model"
)

/*AddBoard to Workplace*/
func (storage *Storage) AddBoard(board model.Board) {
	storage.Workplace.Lanes = append(storage.Workplace.Lanes, board)
	storage.Save()
}

/*AddCard to board of Workplace*/
func (storage *Storage) AddCard(board string, card model.BoardItem) {
	number, err := strconv.ParseInt(board, 10, 32)
	if err == nil {
		board := int(number)
		if len(storage.Workplace.Lanes) >= board {
			storage.Workplace.Lanes[board-1].Items = append(storage.Workplace.Lanes[board-1].Items, card)
			storage.Save()
		}
	} else {
		fmt.Println(err)
	}
}
