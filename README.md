# goyoubbs

## 快速使用

```bash
CGO_ENABLED=0 GOOS=linux go build
docker build . -t yiqi:dev

docker stack deploy --compose-file docker-compose.yaml yiqi
```