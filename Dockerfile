# syntax=docker/dockerfile:1

# ---- Build ----
FROM golang:1.27-alpine AS build

WORKDIR /src

# Abhängigkeiten getrennt laden, damit der Layer gecacht wird
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .

# -trimpath:          keine absoluten Build-Pfade im Binary
# -buildvcs=false:    keine Commit-/Repo-Infos im Binary
# -s -w -buildid=:    keine Symbol-/Debug-Tabellen, keine Build-ID
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -buildvcs=false -ldflags="-s -w -buildid=" -o /out/kanbantabs . \
 && mkdir -p /out/home/kanban \
 && echo "kanban:x:10001:10001:kanban:/home/kanban:/sbin/nologin" > /out/passwd \
 && echo "kanban:x:10001:" > /out/group

# ---- Runtime ----
FROM scratch

# os/user liest ohne cgo /etc/passwd, um das Home-Verzeichnis zu finden
COPY --from=build /out/passwd /etc/passwd
COPY --from=build /out/group /etc/group
COPY --from=build --chown=10001:10001 /out/home/kanban /home/kanban
COPY --from=build /out/kanbantabs /kanbantabs

USER 10001:10001
WORKDIR /home/kanban

VOLUME /home/kanban

ENTRYPOINT ["/kanbantabs"]
