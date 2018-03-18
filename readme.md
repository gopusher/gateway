# Comet 接入层

> 支持包括 TCP、WebSocket 等多种通信协议接入，支持包括 Json、Redis 等多种消息协议

# 编译
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```

## 使用方法
```
./chat -c ./comet.ini
```
