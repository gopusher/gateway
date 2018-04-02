# Comet 接入层

> 支持包括 TCP、WebSocket 等多种通信协议接入，支持包括 Json、Redis 等多种消息协议

## 编译
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```

## 使用方法

1. 启动**etcd**服务，可以选用已有的，如果没有，则可以用以下[这个](https://github.com/Gopusher/awesome/blob/master/docker/docker-compose.yml) （包含了redis和etcd，如果不需要redis则去掉redis部分）

   ```
   docker-compose up -d
   ```

2. 配置

   > 启动服务的时候指定 `-c`参数指定配置文件

   ```
   # 服务最大的cpu执行数
   max_proc=2
   # etcd 中 comet service name
   comet_service_name=comet
   # etcd server addr
   etcd_addr=127.0.0.1:2379
   # 推送逻辑服务器(message 服务) addr
   rpc_api_url=http://msg.demo.com/im/index/rpc
   # 和 推送逻辑服务器(message server) rpc服务通信的 user-agent
   rpc_user_agent="CtxImRpc 1.0"

   # 通信协议 可选项 websocket, tcp (目前仅支持websocket，tcp需要后续开发)
   socket_protocol=websocket
   # websockeet 监听地址
   websocket_host=comet.demo.com
   websocket_port=:8900
   # websocket 协议，可选项 ws, wss (如果为 wss 需要设置 wss_cert_pem 和 wss_key_pem)
   websocket_protocol=ws
   # wss_cert_pem=/path/fullchain.pem
   # wss_key_pem=/path/privkey.pem
   # comet服务的 rpc 监听端口
   rpc_addr=10.0.1.131:8901
   # rpc_addr=192.168.31.86:8901
   ```

2. monitor 服务启动

   监管 come 接入层服务状态上线下线情况 同时通知 逻辑层api.

   ```
   ./chat -m=true -c=./comet.ini 
   ```

4. comet 接入层服务启动

   ```
   ./chat -c=./comet.ini
   ```

## 时序图 

> 推送逻辑服务器(message 服务) 是提供路由服务和业务逻辑服务等，需要由业务方自己实现([参考](https://github.com/Gopusher/message)).

![Comet 接入层服务启动时序图](https://raw.githubusercontent.com/Gopusher/comet/master/docs/Comet%E6%8E%A5%E5%85%A5%E5%B1%82%E6%9C%8D%E5%8A%A1%E5%90%AF%E5%8A%A8%E6%97%B6%E5%BA%8F%E5%9B%BE.png)

![client客户端上线下线时序图](https://raw.githubusercontent.com/Gopusher/comet/master/docs/Client%E5%AE%A2%E6%88%B7%E7%AB%AF%E4%B8%8A%E7%BA%BF%E4%B8%8B%E7%BA%BF%E6%97%B6%E5%BA%8F%E5%9B%BE.png)

![消息发送接收时序图](https://raw.githubusercontent.com/Gopusher/comet/master/docs/%E6%B6%88%E6%81%AF%E5%8F%91%E9%80%81%E6%8E%A5%E6%94%B6%E6%97%B6%E5%BA%8F%E5%9B%BE.png)

## 推送逻辑服务器(message server) 开发

* Api 相关

  > message服务需要提供api，供 monitor 服务 和 comet服务 回调

  1. Comet 接入层服务上线回调接口
  2. Comet 接入层服务下线回调接口
  3. Comet 接入层服务校验Client 客户端Token接口
  4. Comet 接入层服务通知message 服务 client 上线接口
  5. Comet 接入层服务通知message 服务 client 下线接口

* Rpc 相关

  > message服务调用 Comet 服务发送消息

  1. 消息发送接口

## TODO 

* 增加对客户端连接超时的控制处理
* 增加信号接管服务的平滑重启和关闭
* 增加 rpc 方法获取服务 status
* 增加服务负载自检，过载则回调路由服务停止对当前机器的服务下发，防止更多的客户端连接
* 完善通信协议文档
* 增加对tcp协议接入的支持
