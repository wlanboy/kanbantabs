package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	model "github.com/wlanboy/kanbantabs/v2/model"
	store "github.com/wlanboy/kanbantabs/v2/store"
)

const (
	minNameLength = 3
	maxNameLength = 20
)

/*ErrUsage marks wrong command line usage; the caller should print the help text*/
var ErrUsage = errors.New("invalid usage")

/*Usage is the help text for the command line*/
const Usage = `Usage: kanbantabs [object verb [arguments]]

  kanbantabs                               show the board
  kanbantabs lane add [name]               add a lane at the end
  kanbantabs lane rename <lane> [name]     rename a lane
  kanbantabs lane delete <lane>            delete a lane and its cards
  kanbantabs card add [lane [name]]        add a card (default: lane 1)
  kanbantabs card move <id> [lane]         move a card to the next or the given lane
  kanbantabs card rename <id> [name]       rename a card
  kanbantabs card severity <id> <level>    set severity: low, medium (*), high (!)
  kanbantabs card delete <id>              delete a card
  kanbantabs help                          show this help

Lanes are numbered from 1 (left to right), cards by the id shown in brackets.
Names without an argument are asked for interactively; they need 3 to 20 characters.
The board is stored in ~/.kanbantabs or in the file named by KANBANTABS_FILE.`

/*Runner executes command line arguments against a storage*/
type Runner struct {
	Storage *store.Storage
	In      *bufio.Reader
	Out     io.Writer
}

/*NewRunner creates a Runner reading answers from in and writing to out*/
func NewRunner(storage *store.Storage, in io.Reader, out io.Writer) *Runner {
	return &Runner{Storage: storage, In: bufio.NewReader(in), Out: out}
}

/*Run executes the arguments (without program name) and prints the board after changes*/
func (r *Runner) Run(args []string) error {
	if len(args) == 0 {
		Show(r.Storage, r.Out)
		return nil
	}
	if len(args) == 1 && isHelp(args[0]) {
		fmt.Fprintln(r.Out, Usage)
		return nil
	}
	if len(args) < 2 {
		return fmt.Errorf("%w: missing verb for %q", ErrUsage, args[0])
	}

	object, verb, params := args[0], args[1], args[2:]
	var err error
	switch object {
	case "lane":
		err = r.lane(verb, params)
	case "card":
		err = r.card(verb, params)
	default:
		err = fmt.Errorf("%w: unknown object %q", ErrUsage, object)
	}
	if err != nil {
		return err
	}
	Show(r.Storage, r.Out)
	return nil
}

func (r *Runner) lane(verb string, params []string) error {
	switch verb {
	case "add":
		name, err := r.name("lane name", params)
		if err != nil {
			return err
		}
		return r.Storage.AddLane(name)
	case "rename":
		number, err := r.laneNumber(params)
		if err != nil {
			return err
		}
		name, err := r.name("new lane name", params[1:])
		if err != nil {
			return err
		}
		return r.Storage.RenameLane(number, name)
	case "delete":
		if err := expectArgs(params, 1); err != nil {
			return err
		}
		number, err := r.laneNumber(params)
		if err != nil {
			return err
		}
		lane, err := r.Storage.Lane(number)
		if err != nil {
			return err
		}
		if count := len(lane.Items); count > 0 {
			question := fmt.Sprintf("Lane %q still has %d card(s). Delete anyway? [y/N]", lane.Name, count)
			if !r.confirm(question) {
				return errors.New("aborted, nothing deleted")
			}
		}
		return r.Storage.DeleteLane(number)
	default:
		return fmt.Errorf("%w: unknown verb %q for lane", ErrUsage, verb)
	}
}

func (r *Runner) card(verb string, params []string) error {
	switch verb {
	case "add":
		number := 1
		if len(params) > 0 {
			var err error
			if number, err = r.laneNumber(params); err != nil {
				return err
			}
			params = params[1:]
		}
		// check the lane before asking for a name
		if _, err := r.Storage.Lane(number); err != nil {
			return err
		}
		name, err := r.name("card name", params)
		if err != nil {
			return err
		}
		_, err = r.Storage.AddCard(number, name)
		return err
	case "move":
		id, err := cardID(params)
		if err != nil {
			return err
		}
		switch len(params) {
		case 1:
			return r.Storage.MoveCard(id)
		case 2:
			number, err := r.laneNumber(params[1:])
			if err != nil {
				return err
			}
			return r.Storage.MoveCardTo(id, number)
		default:
			return fmt.Errorf("%w: too many arguments", ErrUsage)
		}
	case "rename":
		id, err := cardID(params)
		if err != nil {
			return err
		}
		name, err := r.name("new card name", params[1:])
		if err != nil {
			return err
		}
		return r.Storage.RenameCard(id, name)
	case "severity":
		if err := expectArgs(params, 2); err != nil {
			return err
		}
		id, err := cardID(params)
		if err != nil {
			return err
		}
		severity, err := model.ParseSeverity(params[1])
		if err != nil {
			return err
		}
		return r.Storage.SetSeverity(id, severity)
	case "delete":
		if err := expectArgs(params, 1); err != nil {
			return err
		}
		id, err := cardID(params)
		if err != nil {
			return err
		}
		return r.Storage.DeleteCard(id)
	default:
		return fmt.Errorf("%w: unknown verb %q for card", ErrUsage, verb)
	}
}

// name takes the name from the remaining arguments or asks for it.
func (r *Runner) name(question string, params []string) (string, error) {
	var name string
	if len(params) > 0 {
		name = strings.Join(params, " ")
	} else {
		answer, err := r.ask(fmt.Sprintf("Enter %s:", question))
		if err != nil {
			return "", err
		}
		name = answer
	}
	return normalizeName(name)
}

func (r *Runner) confirm(question string) bool {
	answer, err := r.ask(question)
	if err != nil {
		return false
	}
	answer = strings.ToLower(answer)
	return answer == "y" || answer == "yes"
}

func (r *Runner) ask(question string) (string, error) {
	fmt.Fprintln(r.Out, question)
	line, err := r.In.ReadString('\n')
	if err != nil && (!errors.Is(err, io.EOF) || line == "") {
		return "", errors.New("no input")
	}
	return strings.TrimSpace(line), nil
}

// normalizeName trims the name and cuts it to maxNameLength characters (not bytes).
func normalizeName(name string) (string, error) {
	name = strings.Join(strings.Fields(name), " ")
	if utf8.RuneCountInString(name) < minNameLength {
		return "", fmt.Errorf("name too short, at least %d characters", minNameLength)
	}
	if runes := []rune(name); len(runes) > maxNameLength {
		name = strings.TrimSpace(string(runes[:maxNameLength]))
	}
	return name, nil
}

func (r *Runner) laneNumber(params []string) (int, error) {
	if len(params) == 0 {
		return 0, fmt.Errorf("%w: missing lane number", ErrUsage)
	}
	number, err := strconv.Atoi(params[0])
	if err != nil {
		return 0, fmt.Errorf("lane must be a number, got %q", params[0])
	}
	return number, nil
}

func cardID(params []string) (int32, error) {
	if len(params) == 0 {
		return 0, fmt.Errorf("%w: missing card id", ErrUsage)
	}
	id, err := strconv.ParseInt(params[0], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("card id must be a number, got %q", params[0])
	}
	return int32(id), nil
}

func expectArgs(params []string, count int) error {
	if len(params) != count {
		return fmt.Errorf("%w: expected %d argument(s), got %d", ErrUsage, count, len(params))
	}
	return nil
}

func isHelp(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}
