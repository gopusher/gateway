# Gopusher Comet

Gopusher 是一个支持分布式部署的通用长连接接入层服务，接管客户端连接。Gopusher Comet is a access layer service that handling all client persistent connections with distributed cluster deployment.

你可以很容易的使用 **http api** 来构建实时聊天，通知推送应用。You can use **http api** to develop a instant messaging application or a push notification application easily.

> demo: [https://chat.yadou.net](https://chat.yadou.net)
>
> 这是一个用php开发的聊天应用([源码](https://github.com/Gopusher/laravel-chat))，采用comet作为长连接接入层，采用php开发聊天的路由和逻辑层部分。
>
> This is a chat application developed in php ([code souce](https://github.com/Gopusher/laravel-chat)), using comet service to handle all client persistent connections, and using php language develop the chat routing and  logical layer.

## 开发指南 Develop Guide

* [Wiki](https://github.com/Gopusher/comet/wiki): 包含所有的开发文档和API, All Develop Guide And API Wiki.
* QQ交流群: 818628641

## 特性

* 简单通用
* 多协议支持，websocket 已经支持，tcp 在开发中
* 集群支持
* 开发者友好，采用http api的方式进行rpc调用，便于不同语言的接入开发（你不需要学习golang，只需要运行comet。）
* 支持 ***docker-compose方式运行*** ([docker-compose 运行](https://github.com/Gopusher/awesome/tree/master/docker))

## Features

* light weight
* multi-protocol support, websocket is already supported, tcp is coming soon
* cluster support
* developer friendly, rpc call using http api to make develop with any program languages easily( You don't need to learn golang, just run comet. )
* Support running with docker  ([running with docker](https://github.com/Gopusher/awesome/tree/master/docker))

## 安装 Installation

### Installing Go

[https://golang.org/doc/install](https://golang.org/doc/install)

### 下载 Download

下载项目源码。download comet souce code.

### 依赖 Dependencies

```
go get github.com/coreos/etcd/clientv3
go get github.com/gorilla/websocket
go get github.com/fatih/color
go get gopkg.in/ini.v1
```

### 编译 Build

```
go build -o comet main.go
```

mac上编译debian版本, build debian bin on mac
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o comet-for-debian main.go
```

## 运行 Run

### 运行etcd Run etcd

[run with docker](https://github.com/Gopusher/awesome/blob/master/docker/docker-compose.yml)

当然你可以选择你喜欢的方式运行etcd.  Of course you can run etcd as the way you like.

### 配置 Configuration

> Edit `comet.ini`

```
# 最大的cpu执行数
MAX_PROC=1
# etcd 中 comet service name
COMET_SERVICE_NAME=comet
# etcd server addr
ETCD_ADDR=127.0.0.1:2379
# api addr
RPC_API_URL=http://www.chat.com/im/index/rpc
RPC_USER_AGENT="CtxImRpc 1.0"

# 通信协议 协议，可选项 tcp(tcp需要后续开发), ws, wss (如果为 wss 需要设置 WSS_CERT_PEM 和 WSS_KEY_PEM)
SOCKET_PROTOCOL=ws
# websockeet 监听端口
WEBSOCKET_PORT=8900
# WSS_CERT_PEM=
# WSS_KEY_PEM=
# comet rpc 配置
COMET_RPC_ADDR=192.168.3.142
COMET_RPC_PORT=8901
COMET_RPC_TOKEN=token
```

### 运行 Run
1. Run comet monitor

```
./comet -c comet.ini -m
```

2. Run comet

```
./comet -c comet.ini
```
到现在为止，你已经可以使用comet了，并采用你喜欢的语言进行接入开发你的长连接应用了。So far, you can already use comet service and develop your persistent connections application with your favorite program language.

### 集群配置 Cluster configuration 

如果你需要采用集群的方式运行，你可以采用nginx等来做负载均衡。If you need to run comet cluster, you can use nginx, etc. for load balancing.

```
upstream websocket {
    server 192.168.3.165:8900 weight=1;
    server 192.168.3.165:8902 weight=1;
}

server {
    listen 8910;

    server_name www.chat.com$;

    location / {
        proxy_pass http://websocket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

你也可以采用服务器下发comet ip port的方式来进行负载均衡，you can also use the method of sending the comet ip port to load balance.

### Todo
* 所有的配置都要在log中体现 方便错误信息排查
* 加入集群之前先等待api服务和ws服务启动后

