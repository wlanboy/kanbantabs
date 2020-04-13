FROM busybox:1.31

COPY ./kanbantabs /home/
EXPOSE 8000

CMD ["/home/kanbantabs"]
