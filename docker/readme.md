# Run with Docker

## 构建镜像 Build

```
docker build -t comet .
```

## 配置 Configuration

> Edit `comet.ini`

```
# 最大的cpu执行数
max_proc=1
# etcd 中 comet service name
comet_service_name=comet
# etcd server addr
etcd_addr=192.168.3.165:2379
# api addr
rpc_api_url=http://192.168.3.165/im/index/rpc
rpc_user_agent="CtxImRpc 1.0"

# 通信协议 协议，可选项 tcp(tcp需要后续开发), ws, wss (如果为 wss 需要设置 wss_cert_pem 和 wss_key_pem)
socket_protocol=ws
# websockeet 监听端口
websocket_port=8900
# wss_cert_pem=
# wss_key_pem=
# comet rpc 配置
comet_rpc_addr=192.168.3.165
comet_rpc_port=8901
comet_rpc_token=token
```

## 运行 Run

* 运行etcd Run Etcd

[run with docker](https://github.com/Gopusher/awesome/blob/master/docker/docker-compose.yml)

当然你可以选择你喜欢的方式运行etcd. Of course you can run etcd as the way you like.

* Run comet monitor

```
docker run --rm -it -v $(pwd)/comet.ini:/data/comet.ini comet -c /data/comet.ini -m
```

* Run comet

```
docker run --rm -it -v $(pwd)/comet.ini:/data/comet.ini -p 8900:8900 -p 8901:8901 comet -c /data/comet.ini
```

