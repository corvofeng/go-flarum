# go-flarum


https://github.com/yannisme/flarum-oxo-theme
https://discuss.flarum.org.cn/d/2775

## 本地快速启动

```bash
yarn build
go run main.go
```

```bash
git submodule foreach 'git pull'

git submodule foreach '[[ "$path" =~ (view/extensions/(analytics|auth-github|diff|flarum-pipetables))|(view/locale/zh) ]]  || git checkout main '

git submodule foreach '[[ "$path" =~ view/extensions/(analytics|auth-github|diff|flarum-pipetables) ]]  || echo '
git submodule foreach 'if [[ "$path" =~ view/extensions/(analytics|auth-github|diff|flarum-pipetables) ]]; then git checkout main; fi;'



git submodule foreach '
    if [[ "$path" =~ view/extensions/analytics ]]; then
        git fetch origin && git checkout 1.0.0
    fi

    if [[ "$path" =~ view/extensions/auth-github ]]; then
        git fetch origin && git checkout v0.1.0-beta.13
    fi
'

git submodule foreach 'if [[ "$path" =~ view/extensions/analytics ]]; then git fetch origin && git checkout v1.0.0; fi;'
git submodule foreach 'if [[ "$path" =~ view/extensions/auth-github ]]; then git fetch origin && git checkout v0.1.0-beta.13; fi;'
git submodule foreach 'if [[ "$path" =~ view/extensions/diff ]]; then git fetch origin && git checkout 1.1.1; fi;'
git submodule foreach 'if [[ "$path" =~ view/extensions/(flarum-pipetables) ]]; then git fetch origin && git checkout v2.0; fi;'

```

## 本地使用docker构建

```bash
yarn build
CGO_ENABLED=0 GOOS=linux go build -o go-flarum
docker build . -t go-flarum:dev

docker stack deploy --compose-file docker-compose.yaml go-flarum
```
