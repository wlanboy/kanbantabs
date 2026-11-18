package application

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/joho/godotenv"
	comms "github.com/wlanboy/kanbantabs/v2/commands"
	store "github.com/wlanboy/kanbantabs/v2/store"
)

/*FileEnv names the environment variable that overrides the storage file*/
const FileEnv = "KANBANTABS_FILE"

/*Kanban wires storage and command line together*/
type Kanban struct {
	Storage *store.Storage
}

/*Initialize loads .env (optional) and the board from disk*/
func (kanban *Kanban) Initialize() error {
	// a missing .env file is fine, it only provides optional overrides
	_ = godotenv.Load()

	filename, err := storageFile()
	if err != nil {
		return err
	}
	storage := &store.Storage{Filename: filename}
	if err := storage.Load(); err != nil {
		return err
	}
	kanban.Storage = storage
	return nil
}

/*Run executes the command line arguments (without program name)*/
func (kanban *Kanban) Run(args []string) error {
	return comms.NewRunner(kanban.Storage, os.Stdin, os.Stdout).Run(args)
}

func storageFile() (string, error) {
	if file := os.Getenv(FileEnv); file != "" {
		return file, nil
	}
	if usr, err := user.Current(); err == nil && usr.HomeDir != "" {
		return filepath.Join(usr.HomeDir, ".kanbantabs"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("cannot find home directory, set " + FileEnv)
	}
	return filepath.Join(home, ".kanbantabs"), nil
}

/*PrintError writes err to stderr, with the help text for usage errors*/
func PrintError(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	if errors.Is(err, comms.ErrUsage) {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, comms.Usage)
	}
}
