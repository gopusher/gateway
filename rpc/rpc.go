package rpc

import (
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
	"reflect"
	"bytes"
	"io/ioutil"
	"github.com/fatih/color"
)

type Client struct {
	rpcApiUrl string
	rpcUserAgent string
}

func NewClient(rpcApiUrl string, rpcUserAgent string) *Client {
	return &Client{
		rpcApiUrl: rpcApiUrl,
		rpcUserAgent: rpcUserAgent,
	}
}

func (c Client) send(class string, method string, args ...interface{}) (string, error) {
	type RpcBody struct {
		ClassName	string	`json:"class"`
		MethodName 	string	`json:"method"`
		Args		[]interface{}	`json:"args"`
	}
	body, err := json.Marshal(&RpcBody{
		ClassName: class,
		MethodName: method,
		Args: args,
	})

	req, err := http.NewRequest("POST", c.rpcApiUrl, bytes.NewBuffer(body))
	if err != nil {
		color.Red("rpc 请求失败1 :" + err.Error())
		return "", errors.New("请求失败:" + err.Error())
	}
	req.Header.Set("User-Agent", c.rpcUserAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("rpc 请求失败2 :" + err.Error())
		return "", errors.New("请求失败:" + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		ret, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			color.Red("rpc 请求失败3，StatusCode 200:" + err.Error())
			return "", err
		}
		return string(ret), nil
	}

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red(fmt.Sprintf("rpc 请求异常4, StatusCode %d: %v", resp.StatusCode, err.Error()))
		return "", errors.New(fmt.Sprintf("请求异常%d: %v", resp.StatusCode, err.Error()))
	}

	color.Red(fmt.Sprintf("rpc 请求异常4, StatusCode %d: %s", resp.StatusCode, string(ret)))
	return "", errors.New(string(ret))
}

func (c Client) SuccessRpc(class string, method string, args ...interface{}) (string, error) {
	ret, err := c.send(class, method, args...)
	if err != nil {
		return "", err
	}

	type RetInfo struct {
		Code int `json"code"`
		Data interface{} `json"data"`
		Error string `json"error"`
	}

	var retInfo RetInfo
	if err := json.Unmarshal([]byte(ret), &retInfo); err != nil {
		color.Red("rpc 响应异常, 不能解析 %v %v", ret, reflect.TypeOf(ret))
		return "", errors.New("rpc 响应消息体异常, 不能解析")
	}

	if retInfo.Code != 0 {
		return "", errors.New(retInfo.Error)
	}

	return ret, nil
}
