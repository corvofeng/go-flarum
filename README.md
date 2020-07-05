# zoe

## 本地快速启动

```bash
yarn build
go run main.go
```

## 本地使用docker构建

```bash
yarn build
CGO_ENABLED=0 GOOS=linux go build -o zoe
docker build . -t zoe:dev

docker stack deploy --compose-file docker-compose.yaml zoe
```
