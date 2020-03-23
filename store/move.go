package store

import (
	"strconv"
)

/*MoveCard from board to next board of Workplace*/
func (storage *Storage) MoveCard(cardnumber string) {
	number, err := strconv.ParseInt(cardnumber, 10, 32)
	if err == nil {
		number := int32(number)
		found := false
		foundBoard := 0
		foundCard := 0
		lastboard := len(storage.Workplace.Lanes) - 1

		for _, board := range storage.Workplace.Lanes {

			counter := 0
			cards := board.Items

			for _, card := range cards {
				if number == card.ID {
					found = true
					foundCard = counter
					break
				}
				counter++
			}

			if found {
				cardToMove := cards[foundCard]

				cards[foundCard] = cards[len(cards)-1]
				cards = cards[:len(cards)-1]

				storage.Workplace.Lanes[foundBoard].Items = cards

				if foundBoard < lastboard {
					newBoard := foundBoard + 1
					storage.Workplace.Lanes[newBoard].Items = append(storage.Workplace.Lanes[newBoard].Items, cardToMove)
				}

				storage.Save()

				break
			}
			foundBoard++
		} //for _, board := range lanes
	}
}
