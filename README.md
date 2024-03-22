
```
docker build -t musicbot-img .

docker run -d --name musicbot --restart always -v $PWD/config.yaml:/app/config.yaml -it musicbot-img
```

