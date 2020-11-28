![Go](https://github.com/wlanboy/kanbantabs/workflows/Go/badge.svg?branch=master) ![Docker build and publish image](https://github.com/wlanboy/kanbantabs/workflows/Docker%20build%20and%20publish%20image/badge.svg)

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
* CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v .

# go lang build for other archs
* GOOS=linux GOARCH=386 go build (386 needed for busybox)
* GOOS=linux GOARCH=arm GOARM=6 go build (Raspberry Pi build)
* GOOS=linux GOARCH=arm64 go build (Odroid C2 build)

# docker build
docker build -t kanbantabs:latest . --build-arg BIN_FILE=./kanbantabs

# image size
* 777.72 KB (!)

# docker hub
* https://hub.docker.com/r/wlanboy/kanbantabs

# usage with docker
- alias kanban="docker run -it --rm -v /yourlocalstorage:/home/kanban wlanboy/kanbantabs"
- kanban
- kanban lane add
- kanban lane add
- kanban card add
- kanban card move 1

# usage with binary
* ./kanbantabs (show kanban board)
* ./kanbantabs lane add (add new lane, you will be asked for it's name)
* ./kanbantabs lane delete 1 (delete first lane)
* ./kanbantabs card add (add a new card to the first lane, you will be asked for it's name)
* ./kanbantabs card add 2 (add a new card to the second lane, you will be asked for it's name)
* ./kanbantabs card move 8 (move card with number 8 to the next lane - will disappear after last lane) 
* ./kanbantabs card delete 8 (delete card with number 8) 
