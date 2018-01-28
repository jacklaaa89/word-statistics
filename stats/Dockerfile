FROM alpine

ENV PORT=8080 FILE=/var/data/data.txt
COPY . /var/app/.
WORKDIR /var/app/
RUN mkdir -p /var/data && touch /var/data/data.txt

ENTRYPOINT ./stats-linux-amd64 -f $FILE -p $PORT
