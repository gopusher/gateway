package main

import (
	"net/rpc/jsonrpc"
	"log"
	"fmt"
	"encoding/json"
	"os"
)

type Message struct {
	To   	[]string	`json:"to"`	//消息接受者
	Msg 	string		`json:"msg"` //为一个json，里边包含 type 消息类型
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("消息发送参数格式: go run ./example.go msg ...to")
	}

	message, _ := json.Marshal(&Message{
		To: os.Args[2:],
		Msg: os.Args[1],
	})

	//连接远程rpc服务
	//这里使用jsonrpc.Dial
	//todo 这里的 ip 要注意
	rpc, err := jsonrpc.Dial("tcp", "127.0.0.1:8901")
	if err != nil {
		log.Fatal(err)
	}
	//response 为 json 字符串
	var response string
	//调用远程方法
	//注意第三个参数是指针类型

	err2 := rpc.Call("Server.SendToConnections", string(message), &response)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(response)
}
