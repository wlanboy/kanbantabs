package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	model "github.com/wlanboy/kanbantabs/v2/model"
)

/*Storage holds the workplace and the file it is persisted to*/
type Storage struct {
	Workplace *model.Workplace
	Filename  string
}

/*
Load reads the workplace from disk. A missing file yields an empty workplace;
an unreadable or corrupt file is an error, so it never gets overwritten.
*/
func (storage *Storage) Load() error {
	file, err := os.ReadFile(storage.Filename)
	if errors.Is(err, os.ErrNotExist) {
		storage.Workplace = &model.Workplace{Name: "Kanban", NextID: 1}
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", storage.Filename, err)
	}

	workplace := model.Workplace{}
	if err := json.Unmarshal(file, &workplace); err != nil {
		return fmt.Errorf("parse %s: %w (file left untouched)", storage.Filename, err)
	}
	ensureNextID(&workplace)
	storage.Workplace = &workplace
	return nil
}

/*Save writes the workplace atomically: temp file in the same directory, then rename*/
func (storage *Storage) Save() error {
	data, err := json.MarshalIndent(storage.Workplace, "", " ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(storage.Filename), filepath.Base(storage.Filename)+".tmp*")
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}
	defer os.Remove(tmp.Name()) // no-op after a successful rename

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("save: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("save: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	if err := os.Rename(tmp.Name(), storage.Filename); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

// ensureNextID keeps ids unique even if NextID in the file is missing or too low.
func ensureNextID(workplace *model.Workplace) {
	if workplace.NextID < 1 {
		workplace.NextID = 1
	}
	for _, lane := range workplace.Lanes {
		for _, card := range lane.Items {
			if card.ID >= workplace.NextID {
				workplace.NextID = card.ID + 1
			}
		}
	}
}
