# 可以指定依赖的node镜像的版本 node:<version>，如果不指定，就会是最新的
FROM node:14.5.0-alpine3.12 as build-static

# 创建工作目录，对应的是应用代码存放在容器内的路径
WORKDIR /usr/src/app

# 把 package.json，package-lock.json(npm@5+) 或 yarn.lock 复制到工作目录(相对路径)
COPY package.json *.lock ./

# 只安装dependencies依赖
# node镜像自带yarn
# RUN yarn --only=prod --registry=https://registry.npm.taobao.org
RUN yarn --only=prod

COPY view webpack.config.js ./


RUN yarn build

# Golang编译阶段
FROM 1.14.4-alpine3.12 as build-backend
# All these steps will be cached
WORKDIR /home/zoe

# 国内用户可能需要设置 go proxy
# RUN go env -w GOPROXY=https://goproxy.cn,direct

# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/hello

FROM alpine:3.7

WORKDIR /home/zoe

COPY ./goyoubbs $PROJDIR/goyoubbs
COPY ./config/config.yaml $PROJDIR/config.yml
COPY ./static $PROJDIR/static
COPY ./view $PROJDIR/view

EXPOSE 8082
CMD ["/home/goyoubbs/goyoubbs", "-config", "/home/goyoubbs/config.yml"]
