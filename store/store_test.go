package store

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	model "github.com/wlanboy/kanbantabs/v2/model"
)

func newStorage(t *testing.T) *Storage {
	t.Helper()
	storage := &Storage{Filename: filepath.Join(t.TempDir(), ".kanbantabs")}
	if err := storage.Load(); err != nil {
		t.Fatal(err)
	}
	return storage
}

func laneNames(storage *Storage) []string {
	var names []string
	for _, lane := range storage.Workplace.Lanes {
		names = append(names, lane.Name)
	}
	return names
}

func cardIDs(lane model.Board) []int32 {
	var ids []int32
	for _, card := range lane.Items {
		ids = append(ids, card.ID)
	}
	return ids
}

func mustAdd(t *testing.T, storage *Storage, lane int, name string) int32 {
	t.Helper()
	card, err := storage.AddCard(lane, name)
	if err != nil {
		t.Fatal(err)
	}
	return card.ID
}

func TestLoadMissingFileGivesEmptyWorkplace(t *testing.T) {
	storage := newStorage(t)
	if storage.Workplace.NextID != 1 || len(storage.Workplace.Lanes) != 0 {
		t.Fatalf("unexpected workplace: %+v", storage.Workplace)
	}
	if _, err := os.Stat(storage.Filename); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load must not create the file, stat: %v", err)
	}
}

func TestLoadCorruptFileFailsAndKeepsFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), ".kanbantabs")
	corrupt := []byte(`{"Name": "Kanban", "Lanes": [`)
	if err := os.WriteFile(filename, corrupt, 0600); err != nil {
		t.Fatal(err)
	}

	storage := &Storage{Filename: filename}
	if err := storage.Load(); err == nil {
		t.Fatal("expected error for corrupt file")
	}
	content, _ := os.ReadFile(filename)
	if string(content) != string(corrupt) {
		t.Fatalf("corrupt file was modified: %q", content)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	storage := newStorage(t)
	storage.AddLane("dev")
	id := mustAdd(t, storage, 1, "first")
	storage.SetSeverity(id, model.HIGH)

	loaded := &Storage{Filename: storage.Filename}
	if err := loaded.Load(); err != nil {
		t.Fatal(err)
	}
	card := loaded.Workplace.Lanes[0].Items[0]
	if card.Name != "first" || card.Severity != model.HIGH || loaded.Workplace.NextID != 2 {
		t.Fatalf("unexpected state after reload: %+v", loaded.Workplace)
	}

	entries, _ := os.ReadDir(filepath.Dir(storage.Filename))
	if len(entries) != 1 {
		t.Fatalf("temp files left behind: %v", entries)
	}
}

func TestLoadRepairsNextID(t *testing.T) {
	filename := filepath.Join(t.TempDir(), ".kanbantabs")
	data := `{"Name":"Kanban","NextID":0,"Lanes":[{"Name":"dev","Items":[{"ID":7,"Name":"x"}]}]}`
	os.WriteFile(filename, []byte(data), 0600)

	storage := &Storage{Filename: filename}
	if err := storage.Load(); err != nil {
		t.Fatal(err)
	}
	if storage.Workplace.NextID != 8 {
		t.Fatalf("NextID = %d, want 8", storage.Workplace.NextID)
	}
}

func TestDeleteLaneKeepsOrder(t *testing.T) {
	storage := newStorage(t)
	for _, name := range []string{"dev", "test", "prod"} {
		storage.AddLane(name)
	}
	if err := storage.DeleteLane(1); err != nil {
		t.Fatal(err)
	}
	if got := laneNames(storage); !slices.Equal(got, []string{"test", "prod"}) {
		t.Fatalf("lanes = %v", got)
	}
	if err := storage.DeleteLane(3); !errors.Is(err, ErrLaneNotFound) {
		t.Fatalf("expected ErrLaneNotFound, got %v", err)
	}
}

func TestAddCardInvalidLaneKeepsNextID(t *testing.T) {
	storage := newStorage(t)
	if _, err := storage.AddCard(1, "card"); !errors.Is(err, ErrLaneNotFound) {
		t.Fatalf("expected ErrLaneNotFound, got %v", err)
	}
	if storage.Workplace.NextID != 1 {
		t.Fatalf("NextID changed to %d", storage.Workplace.NextID)
	}
}

func TestDeleteCardKeepsOrder(t *testing.T) {
	storage := newStorage(t)
	storage.AddLane("dev")
	first := mustAdd(t, storage, 1, "a")
	mustAdd(t, storage, 1, "b")
	mustAdd(t, storage, 1, "c")

	if err := storage.DeleteCard(first); err != nil {
		t.Fatal(err)
	}
	if got := cardIDs(storage.Workplace.Lanes[0]); !slices.Equal(got, []int32{2, 3}) {
		t.Fatalf("cards = %v", got)
	}
	if err := storage.DeleteCard(99); !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("expected ErrCardNotFound, got %v", err)
	}
}

func TestMoveCard(t *testing.T) {
	storage := newStorage(t)
	storage.AddLane("dev")
	storage.AddLane("prod")
	a := mustAdd(t, storage, 1, "a")
	b := mustAdd(t, storage, 1, "b")
	c := mustAdd(t, storage, 1, "c")

	if err := storage.MoveCard(a); err != nil {
		t.Fatal(err)
	}
	if got := cardIDs(storage.Workplace.Lanes[0]); !slices.Equal(got, []int32{b, c}) {
		t.Fatalf("dev = %v", got)
	}
	if got := cardIDs(storage.Workplace.Lanes[1]); !slices.Equal(got, []int32{a}) {
		t.Fatalf("prod = %v", got)
	}

	// last lane: card stays
	if err := storage.MoveCard(a); !errors.Is(err, ErrLastLane) {
		t.Fatalf("expected ErrLastLane, got %v", err)
	}
	if got := cardIDs(storage.Workplace.Lanes[1]); !slices.Equal(got, []int32{a}) {
		t.Fatalf("card left last lane: %v", got)
	}

	// back to a given lane
	if err := storage.MoveCardTo(a, 1); err != nil {
		t.Fatal(err)
	}
	if got := cardIDs(storage.Workplace.Lanes[0]); !slices.Equal(got, []int32{b, c, a}) {
		t.Fatalf("dev = %v", got)
	}
	if err := storage.MoveCardTo(a, 5); !errors.Is(err, ErrLaneNotFound) {
		t.Fatalf("expected ErrLaneNotFound, got %v", err)
	}
}

func TestRename(t *testing.T) {
	storage := newStorage(t)
	storage.AddLane("dev")
	id := mustAdd(t, storage, 1, "old")

	if err := storage.RenameLane(1, "develop"); err != nil {
		t.Fatal(err)
	}
	if err := storage.RenameCard(id, "new"); err != nil {
		t.Fatal(err)
	}
	lane := storage.Workplace.Lanes[0]
	if lane.Name != "develop" || lane.Items[0].Name != "new" {
		t.Fatalf("unexpected lane: %+v", lane)
	}
}
