# Gopusher

Gopusher 是一个支持分布式部署的通用长连接接入层服务，接管客户端连接。Gopusher Comet is a access layer service that handling all client persistent connections with distributed cluster deployment.

你可以很容易的使用 **http api** 来构建实时聊天，通知推送等应用。You can use **http api** to develop a instant messaging application or a push notification application easily.

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

### 编译 Build

```
CGO_ENABLED=0 go build -ldflags '-s -w' -o gopusher main.go
```

mac上编译linux版本, build linux bin file on mac
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o gopusher-for-linux main.go
```

## 运行 Run

### 配置 Configuration

> Edit `vim config.yaml`

```
app_name: gateway

app_debug: true

logging:
    level: warn
    log_file: ""

api_server:
    address: 0.0.0.0:8080
    token: "token"

redis:
    default:
        addr: 127.0.0.1:6379
        password: ""
        db: 0

dingtalk:
    default:
        token: "xx"
        secret: "yy"

node: gateway-01

server.websocket:
    app_key: "app_key_1"
    app_secret: "app_secret_1"
    address: :8900
    ssl: false
    ssl_cert_file: ""
    ssl_key_file: ""
    client_id_alias: "client_id"
    token_alias: "token"
    time_alias: "time"
    time_window: 5
```

### 运行 Run
1. Run Gopusher

```
./gopusher start -c config.yaml
```
到现在为止，你已经可以使用gopusher了，并采用你喜欢的语言进行接入开发你的长连接应用了。So far, you can already use gopusher service and develop your persistent connections application with your favorite program language.

## Docker Usage

1. Build docker image

   ```
   docker build -t gopusher .
   ```

1. Run

  ```
  docker run --rm -it -v $(pwd)/config.yaml:/app/config.yaml --net=host gopusher gopusher start
  ```

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

你也可以采用服务器下发 gopusher ip port 的方式来进行负载均衡，you can also use the method of sending the gopusher ip port to load balance.

### SEO Strings
Words 'chat' and 'instant messaging' in Chinese, Russian, Persian and a few other languages.

* 聊天室 即時通訊
* чат мессенджер
* インスタントメッセージ
* 인스턴트 메신저
* پیام‌رسانی فوری گپ
* تراسل فوري
* Nhắn tin tức thời
* Mensageiro instantâneo
* Pesan instan

## Todo
* 增加 prometheus 的支持，如 当前所有在线连接数 等指标
* 增加etcd服务注册
* 增加服务运行状态接口：协程数，内存占用，当前连接数，服务运行开始时间，已经处理消息数等

