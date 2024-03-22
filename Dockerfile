FROM golang:1.14-alpine3.12

RUN apk update && apk add git ffmpeg ca-certificates && update-ca-certificates

WORKDIR /app

COPY . .

RUN go build -o go-discord-music-bot

RUN mkdir -p /app/queue

ENTRYPOINT ["/app/go-discord-music-bot"]

