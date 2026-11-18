package store

import (
	"errors"
	"fmt"
	"slices"

	model "github.com/wlanboy/kanbantabs/v2/model"
)

/*ErrCardNotFound is returned for unknown card ids*/
var ErrCardNotFound = errors.New("card not found")

/*AddCard appends a new card to the lane with the given 1-based number*/
func (storage *Storage) AddCard(laneNumber int, name string) (model.BoardItem, error) {
	lane, err := storage.Lane(laneNumber)
	if err != nil {
		return model.BoardItem{}, err
	}
	card := model.BoardItem{ID: storage.Workplace.NextID, Name: name}
	lane.Items = append(lane.Items, card)
	storage.Workplace.NextID++
	return card, storage.Save()
}

/*DeleteCard removes a card, keeping the order of the other cards*/
func (storage *Storage) DeleteCard(id int32) error {
	laneIndex, cardIndex, err := storage.findCard(id)
	if err != nil {
		return err
	}
	lane := &storage.Workplace.Lanes[laneIndex]
	lane.Items = slices.Delete(lane.Items, cardIndex, cardIndex+1)
	return storage.Save()
}

/*RenameCard changes the name of a card*/
func (storage *Storage) RenameCard(id int32, name string) error {
	card, err := storage.card(id)
	if err != nil {
		return err
	}
	card.Name = name
	return storage.Save()
}

/*SetSeverity changes the severity of a card*/
func (storage *Storage) SetSeverity(id int32, severity model.SeverityType) error {
	card, err := storage.card(id)
	if err != nil {
		return err
	}
	card.Severity = severity
	return storage.Save()
}

func (storage *Storage) card(id int32) (*model.BoardItem, error) {
	laneIndex, cardIndex, err := storage.findCard(id)
	if err != nil {
		return nil, err
	}
	return &storage.Workplace.Lanes[laneIndex].Items[cardIndex], nil
}

// findCard returns the 0-based lane and card index of the card with the given id.
func (storage *Storage) findCard(id int32) (laneIndex int, cardIndex int, err error) {
	for laneIndex, lane := range storage.Workplace.Lanes {
		for cardIndex, card := range lane.Items {
			if card.ID == id {
				return laneIndex, cardIndex, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("%w: %d", ErrCardNotFound, id)
}
