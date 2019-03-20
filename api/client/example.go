package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"os"
)

type Message struct {
	Connections []string `json:"connections"` //消息接受者
	Msg         string   `json:"msg"`         //为一个json，里边包含 type 消息类型
	Token       string   `json:"token"`       //作为消息发送鉴权
}

type KickMessage struct {
	Connections []string `json:"connections"` //消息接受者
	Token       string   `json:"token"`       //作为消息发送鉴权
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("消息发送参数格式: go run ./example.go msg ...to")
	}

	//连接远程rpc服务
	//这里使用jsonrpc.Dial
	//todo 这里的 ip 要注意
	rpc, err := jsonrpc.Dial("tcp", "192.168.3.165:8901")
	if err != nil {
		log.Fatal(err)
	}
	//response 为 json 字符串
	var response string
	//调用远程方法
	//注意第三个参数是指针类型

	//发送消息
	err2 := rpc.Call("Server.SendToConnections", &Message{
		Connections: os.Args[2:],
		Msg:         os.Args[1],
		Token:       "token",
	}, &response)

	//kick conns
	//err2 := rpc.Call("Server.KickConnections", &KickMessage{
	//	Connections: os.Args[1:],
	//	Token: "token",
	//}, &response)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(response)
}
