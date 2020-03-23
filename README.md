# kanbantabs
golang based console kanban board

# build
* go get -d -v
* go clean
* go build

# run
* go run main.go

# install
* go install

# usage
* ./kanbantabs (show kanban board)
* ./kanbantabs lane add (add new lange, you will be asked for it's name)
* ./kanbantabs lane delete 1 (delete first lane)
* ./kanbantabs card add 1 (add a new card to the first lane, you will be asked for it's name)
* ./kanbantabs card move 8 (move card with number 8 to the nex lane - will disappear after last lane) 
* ./kanbantabs card delete 8 (delete card with number 8) 
