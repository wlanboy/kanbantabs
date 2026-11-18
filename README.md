![Go](https://github.com/wlanboy/kanbantabs/workflows/Go/badge.svg?branch=main) ![Docker build and publish image](https://github.com/wlanboy/kanbantabs/workflows/Docker%20build%20and%20publish%20image/badge.svg)

# kanbantabs
Kanban board for the console, written in Go. kanbantabs keeps the board in a local JSON file.

Example:
```bash
 ./kanbantabs
------------------------------------------------------------------------------
1: dev                   |2: test                  |3: prod                  |
------------------------------------------------------------------------------
[1] kanbantabs           |[3!] ranger              |[4] azure                |
[2] hazelcast            |                         |[6*] gc tuning           |
[5] chromebook           |                         |                         |
```

Lanes are numbered from left to right. Each card shows its id in brackets, followed by `*` for medium or `!` for high severity.

## Usage
```bash
kanbantabs                               # show the board
kanbantabs lane add [name]               # add a lane at the end
kanbantabs lane rename <lane> [name]     # rename a lane
kanbantabs lane delete <lane>            # delete a lane and its cards (asks if the lane has cards)
kanbantabs card add [lane [name]]        # add a card, lane 1 if no lane is given
kanbantabs card move <id> [lane]         # move a card to the next lane or to the given lane
kanbantabs card rename <id> [name]       # rename a card
kanbantabs card severity <id> <level>    # low, medium or high
kanbantabs card delete <id>              # delete a card
kanbantabs help                          # show the help
```

If you leave out the name, kanbantabs asks for it. Names need at least 3 characters, and kanbantabs cuts longer names to 20. A name can have several words: `kanbantabs card add 2 Release vorbereiten`.

A card in the last lane stays there on `card move`. Use `card delete` to remove it.

On errors kanbantabs prints a message to stderr and exits with code 1.

## Storage
kanbantabs stores the board in `~/.kanbantabs`. The environment variable `KANBANTABS_FILE` sets a different file, either directly or through a `.env` file in the current directory.

kanbantabs writes a temp file first and renames it over the old one, so an interrupted save cannot leave a half-written board. If the file cannot be parsed, kanbantabs stops with an error and leaves the file as it is.

## Build and test
```bash
go build
go test ./...
go install
```

Builds for other platforms:
```bash
GOOS=linux GOARCH=arm GOARM=6 go build   # Raspberry Pi
GOOS=linux GOARCH=arm64 go build         # Odroid C2, Raspberry Pi 4/5
```

## Docker
```bash
docker build -t kanbantabs:latest .
```

Image on Docker Hub: https://hub.docker.com/r/wlanboy/kanbantabs

```bash
alias kanban="docker run -it --rm -v /yourlocalstorage:/home/kanban wlanboy/kanbantabs"
kanban
kanban lane add dev
kanban lane add done
kanban card add 1 "first card"
kanban card move 1
```

The host directory must be writable for UID 10001, the user the container runs as.
