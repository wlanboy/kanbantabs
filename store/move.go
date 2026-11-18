package store

import (
	"errors"
	"slices"
)

/*ErrLastLane is returned when a card in the last lane should move to the next lane*/
var ErrLastLane = errors.New("card is already in the last lane")

/*MoveCard moves a card to the next lane. Cards in the last lane stay where they are.*/
func (storage *Storage) MoveCard(id int32) error {
	laneIndex, _, err := storage.findCard(id)
	if err != nil {
		return err
	}
	if laneIndex == len(storage.Workplace.Lanes)-1 {
		return ErrLastLane
	}
	return storage.MoveCardTo(id, laneIndex+2)
}

/*MoveCardTo moves a card to the end of the lane with the given 1-based number*/
func (storage *Storage) MoveCardTo(id int32, laneNumber int) error {
	target, err := storage.Lane(laneNumber)
	if err != nil {
		return err
	}
	laneIndex, cardIndex, err := storage.findCard(id)
	if err != nil {
		return err
	}
	if laneIndex == laneNumber-1 {
		return nil
	}

	source := &storage.Workplace.Lanes[laneIndex]
	card := source.Items[cardIndex]
	source.Items = slices.Delete(source.Items, cardIndex, cardIndex+1)
	target.Items = append(target.Items, card)
	return storage.Save()
}
