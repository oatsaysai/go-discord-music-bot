
```
docker build -t go-discord-music-bot .

docker run -d --name go-discord-music-bot --restart always -v $PWD/config.yaml:/app/config.yaml go-discord-music-bot
```

