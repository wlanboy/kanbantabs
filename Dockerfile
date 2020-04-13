FROM busybox:1.31

ARG BIN_FILE
WORKDIR /
ADD ${BIN_FILE} kanbantabs

ENTRYPOINT ["/kanbantabs"]
