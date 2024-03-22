FROM golang:1.22-alpine3.19

RUN apk update && apk add git ffmpeg ca-certificates && update-ca-certificates

WORKDIR /app

COPY . .

RUN go build -o go-discord-music-bot

RUN mkdir -p /app/queue

ENTRYPOINT ["/app/go-discord-music-bot"]

