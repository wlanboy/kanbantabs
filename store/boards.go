package store

import (
	"errors"
	"fmt"
	"slices"
	"uuid"

	model "github.com/wlanboy/kanbantabs/v2/model"
)

/*ErrLaneNotFound is returned for lane numbers outside 1..len(lanes)*/
var ErrLaneNotFound = errors.New("lane not found")

/*AddLane appends a new lane at the end*/
func (storage *Storage) AddLane(name string) error {
	lane := model.Board{UUID: uuid.NewV4().String(), Name: name}
	storage.Workplace.Lanes = append(storage.Workplace.Lanes, lane)
	return storage.Save()
}

/*Lane returns the lane with the given 1-based number*/
func (storage *Storage) Lane(number int) (*model.Board, error) {
	if number < 1 || number > len(storage.Workplace.Lanes) {
		return nil, fmt.Errorf("%w: %d (there are %d lanes)", ErrLaneNotFound, number, len(storage.Workplace.Lanes))
	}
	return &storage.Workplace.Lanes[number-1], nil
}

/*DeleteLane removes a lane and its cards, keeping the order of the other lanes*/
func (storage *Storage) DeleteLane(number int) error {
	if _, err := storage.Lane(number); err != nil {
		return err
	}
	storage.Workplace.Lanes = slices.Delete(storage.Workplace.Lanes, number-1, number)
	return storage.Save()
}

/*RenameLane changes the name of a lane*/
func (storage *Storage) RenameLane(number int, name string) error {
	lane, err := storage.Lane(number)
	if err != nil {
		return err
	}
	lane.Name = name
	return storage.Save()
}
