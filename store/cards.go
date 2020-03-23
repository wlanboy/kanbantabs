package store

import (
	"fmt"
	"strconv"

	model "../model"
)

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

/*DeleteCard from board of Workplace*/
func (storage *Storage) DeleteCard(cardnumber string) {
	number, err := strconv.ParseInt(cardnumber, 10, 32)
	if err == nil {
		number := int32(number)
		found := false
		foundboard := 0
		foundcard := 0

		for _, board := range storage.Workplace.Lanes {

			counter := 0
			cards := board.Items

			for _, card := range cards {
				if number == card.ID {
					found = true
					foundcard = counter
					break
				}
				counter++
			}

			if found {
				cards[foundcard] = cards[len(cards)-1]
				cards = cards[:len(cards)-1]

				storage.Workplace.Lanes[foundboard].Items = cards
				storage.Save()
				break
			}
			foundboard++
		} //for _, board := range lanes
	}
}
