package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	model "../model"
	store "../store"
	uuid "github.com/satori/go.uuid"
)

/*ExecuteOn for Kanban*/
func ExecuteOn(storage *store.Storage, object string, verb string, param string) {
	switch object {
	case "lane":
		switch verb {
		case "add":
			name := askQuestions("board")
			if len(name) > 2 {
				var board model.Board = model.Board{}
				uid, _ := uuid.NewV4()
				board.UUID = uid.String()
				board.Name = name
				storage.AddBoard(board)
			} else {
				fmt.Print("input to short\n")
			}
		}

	case "card":
		switch verb {
		case "add":
			name := askQuestions("card")
			if len(name) > 2 {
				var card model.BoardItem = model.BoardItem{}
				card.ID = storage.Workplace.NextID
				card.Name = name

				storage.Workplace.NextID++
				storage.AddCard(param, card)
			} else {
				fmt.Print("input to short\n")
			}
		}
	}

}

func askQuestions(question string) string {
	fmt.Printf("Enter %s: \n", question)
	reader := bufio.NewReader(os.Stdin)
	readText, _ := reader.ReadString('\n')
	fmt.Println()

	response := string(readText)
	response = strings.TrimSuffix(response, "\n")
	if len(response) > 20 {
		response = response[0:20]
	}
	return response
}
