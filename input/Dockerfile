FROM alpine

ENV PORT=5555 FILE=/var/data/data.txt
COPY . /var/app/.
WORKDIR /var/app/
RUN mkdir -p /var/data && touch /var/data/data.txt

ENTRYPOINT ./input-linux-amd64 -f $FILE -p $PORT
