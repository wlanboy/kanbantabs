package commands

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	model "github.com/wlanboy/kanbantabs/v2/model"
	store "github.com/wlanboy/kanbantabs/v2/store"
)

func newStorage(t *testing.T, lanes ...string) *store.Storage {
	t.Helper()
	storage := &store.Storage{Filename: filepath.Join(t.TempDir(), ".kanbantabs")}
	if err := storage.Load(); err != nil {
		t.Fatal(err)
	}
	for _, lane := range lanes {
		storage.AddLane(lane)
	}
	return storage
}

func run(storage *store.Storage, input string, args ...string) (string, error) {
	var out bytes.Buffer
	err := NewRunner(storage, strings.NewReader(input), &out).Run(args)
	return out.String(), err
}

func TestNameFromArguments(t *testing.T) {
	storage := newStorage(t, "dev", "test")
	if _, err := run(storage, "", "card", "add", "2", "Release", "vorbereiten"); err != nil {
		t.Fatal(err)
	}
	card := storage.Workplace.Lanes[1].Items[0]
	if card.Name != "Release vorbereiten" {
		t.Fatalf("name = %q", card.Name)
	}
}

func TestNameFromPrompt(t *testing.T) {
	storage := newStorage(t)
	out, err := run(storage, "  dev  \r\n", "lane", "add")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Enter lane name:") || storage.Workplace.Lanes[0].Name != "dev" {
		t.Fatalf("out = %q, lanes = %+v", out, storage.Workplace.Lanes)
	}
}

func TestNameWithoutTrailingNewline(t *testing.T) {
	storage := newStorage(t)
	if _, err := run(storage, "dev", "lane", "add"); err != nil {
		t.Fatal(err)
	}
}

func TestNameErrors(t *testing.T) {
	storage := newStorage(t)
	if _, err := run(storage, "", "lane", "add"); err == nil {
		t.Fatal("expected error for missing input")
	}
	if _, err := run(storage, "", "lane", "add", "ab"); err == nil {
		t.Fatal("expected error for short name")
	}
	if len(storage.Workplace.Lanes) != 0 {
		t.Fatal("lane added despite error")
	}
}

func TestNormalizeNameCutsRunes(t *testing.T) {
	name, err := normalizeName(strings.Repeat("ä", 25))
	if err != nil {
		t.Fatal(err)
	}
	if name != strings.Repeat("ä", 20) {
		t.Fatalf("name = %q", name)
	}
}

func TestCardAddWithoutLanes(t *testing.T) {
	storage := newStorage(t)
	out, err := run(storage, "card\n", "card", "add")
	if !errors.Is(err, store.ErrLaneNotFound) {
		t.Fatalf("expected ErrLaneNotFound, got %v", err)
	}
	if strings.Contains(out, "Enter") {
		t.Fatal("asked for a name although the lane does not exist")
	}
}

func TestLaneDeleteAsksWhenNotEmpty(t *testing.T) {
	storage := newStorage(t, "dev")
	storage.AddCard(1, "card")

	if _, err := run(storage, "n\n", "lane", "delete", "1"); err == nil {
		t.Fatal("expected abort")
	}
	if len(storage.Workplace.Lanes) != 1 {
		t.Fatal("lane deleted without confirmation")
	}
	if _, err := run(storage, "y\n", "lane", "delete", "1"); err != nil {
		t.Fatal(err)
	}
	if len(storage.Workplace.Lanes) != 0 {
		t.Fatal("lane not deleted")
	}
}

func TestCardSeverityAndShow(t *testing.T) {
	storage := newStorage(t, "dev", "prod")
	storage.AddCard(1, "deploy")

	out, err := run(storage, "", "card", "severity", "1", "high")
	if err != nil {
		t.Fatal(err)
	}
	if storage.Workplace.Lanes[0].Items[0].Severity != model.HIGH {
		t.Fatal("severity not set")
	}
	if !strings.Contains(out, "1: dev") || !strings.Contains(out, "[1!] deploy") {
		t.Fatalf("unexpected board:\n%s", out)
	}
}

func TestCardMoveToLane(t *testing.T) {
	storage := newStorage(t, "dev", "test", "prod")
	storage.AddCard(1, "card")
	if _, err := run(storage, "", "card", "move", "1", "3"); err != nil {
		t.Fatal(err)
	}
	if len(storage.Workplace.Lanes[2].Items) != 1 {
		t.Fatal("card not in lane 3")
	}
}

func TestUsageErrors(t *testing.T) {
	storage := newStorage(t, "dev")
	for _, args := range [][]string{
		{"lane"},
		{"board", "add"},
		{"lane", "fly"},
		{"card", "delete"},
		{"card", "delete", "1", "2"},
	} {
		if _, err := run(storage, "", args...); !errors.Is(err, ErrUsage) {
			t.Errorf("%v: expected ErrUsage, got %v", args, err)
		}
	}
	if _, err := run(storage, "", "card", "delete", "x"); err == nil {
		t.Error("expected error for non-numeric id")
	}
}

func TestShowEmpty(t *testing.T) {
	out, err := run(newStorage(t), "")
	if err != nil || !strings.Contains(out, "No lanes yet") {
		t.Fatalf("out = %q, err = %v", out, err)
	}
}
