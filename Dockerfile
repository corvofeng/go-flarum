# 静态资源编译阶段
FROM node:14.16.0-alpine3.12 as build-static

# 创建工作目录，对应的是应用代码存放在容器内的路径
WORKDIR /home/go-flarum
COPY package.json *.lock ./

# 只安装dependencies依赖
# node镜像自带yarn
# ## BOF CLEAN
# # 国内用户可能设置 regietry
ARG registry=https://registry.npm.taobao.org
ARG disturl=https://npm.taobao.org/dist
RUN yarn config set disturl $disturl
RUN yarn config set registry $registry
# ## EOF CLEAN
RUN yarn --only=prod

COPY webpack.config.js ./
COPY view ./view
RUN yarn build

# Golang编译阶段
FROM golang:1.20.5-alpine3.18 as build-backend
# All these steps will be cached
WORKDIR /home/go-flarum

# ## BOF CLEAN
# # 国内用户可能需要设置 go proxy
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN go env -w GOPROXY=https://goproxy.cn,direct
# ## EOF CLEAN

RUN apk update && apk add git

# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

RUN go mod download
# # COPY the source code as the last step
COPY . .

# # Build the binary
RUN GIT_COMMIT=$(git rev-list -1 HEAD) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-a -installsuffix cgo \
	-ldflags "-X main.GitCommit=$GIT_COMMIT" \
	-o go-flarum ./cmd/server/main.go 

RUN GIT_COMMIT=$(git rev-list -1 HEAD) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-a -installsuffix cgo \
	-ldflags "-X main.GitCommit=$GIT_COMMIT" \
	-o go-flarum-migration ./cmd/migration/main.go

# 构建最终镜像
FROM alpine:3.7
WORKDIR /home/go-flarum

COPY ./view view
RUN rm -rf view/extensions view/flarum
COPY ./config/config.yaml-docker config.yml
COPY --from=build-static /home/go-flarum/static webpack/static
COPY --from=build-backend /home/go-flarum/go-flarum /usr/local/bin/go-flarum
COPY --from=build-backend /home/go-flarum/go-flarum-migration /usr/local/bin/go-flarum-migration

EXPOSE 8082
CMD ["/usr/local/bin/go-flarum", "-config", "/home/go-flarum/config.yml"]
