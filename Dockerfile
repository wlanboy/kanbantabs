FROM busybox:1.31

ARG BIN_FILE
ADD ${BIN_FILE} /home/kanbantabs

ENTRYPOINT ["/home/kanbantabs"]
