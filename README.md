![Go](https://github.com/wlanboy/kanbantabs/workflows/Go/badge.svg?branch=master)

# kanbantabs
golang based console kanban board
```
 ./kanbantabs 
------------------------------------------------------------------------------
dev                      |test                     |prod                     |
------------------------------------------------------------------------------
[1] kanbantabs           |[3] ranger               |[4] azure                |
[2] hazelcast            |                         |[6] gc                   |
[5] chromebook           |                         |                         |
```

# build
* go get -d -v
* go clean
* go build

# run
* go run main.go

# install
* go install

# go lang build for docker
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v .

# docker build
docker build -t kanbantabs:latest . --build-arg BIN_FILE=./kanbantabs

# docker run
docker run -it --rm -v /yourlocalstorage:/home/kanban kanbantabs
docker run -it --rm -v /yourlocalstorage:/home/kanban kanbantabs lane add
docker run -it --rm -v /yourlocalstorage:/home/kanban kanbantabs card add 1

# usage
* ./kanbantabs (show kanban board)
* ./kanbantabs lane add (add new lane, you will be asked for it's name)
* ./kanbantabs lane delete 1 (delete first lane)
* ./kanbantabs card add 1 (add a new card to the first lane, you will be asked for it's name)
* ./kanbantabs card move 8 (move card with number 8 to the next lane - will disappear after last lane) 
* ./kanbantabs card delete 8 (delete card with number 8) 
