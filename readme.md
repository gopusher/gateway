# Comet 接入层

> 支持包括 TCP、WebSocket 等多种通信协议接入，支持包括 Json、Redis 等多种消息协议

## 编译
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```

## 使用方法
```
./chat -c ./comet.ini
```

## TODO 
* 增加对客户端连接超时的控制处理
* 增加信号接管服务的平滑重启和关闭
* 增加 rpc 方法获取服务 status
* 增加服务负载自检，过载则回调路由服务停止对当前机器的服务下发，防止更多的客户端连接
