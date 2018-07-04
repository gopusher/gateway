package rpc

import (
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
	"reflect"
	"gopusher/comet/config"
	"bytes"
	"io/ioutil"
	"github.com/fatih/color"
)

type Client struct {
	config *config.Config
}

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
	}
}

func (c Client) Send(class string, method string, args ...interface{}) (string, error) {
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

	apiUrl := c.config.Get("rpc_api_url").String()
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(body))
	if err != nil {
		color.Red("rpc 请求失败1 :" + err.Error())
		return "", errors.New("请求失败:" + err.Error())
	}
	req.Header.Set("User-Agent", c.config.Get("rpc_user_agent").String())
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
			color.Red("rpc 请求失败3，StatusCode 非 200:" + err.Error())
			return "", err
		}
		return string(ret), nil
	}

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red(fmt.Sprintf("rpc 请求异常4, %d: %v", resp.StatusCode, err.Error()))
		return "", errors.New(fmt.Sprintf("请求异常%d: %v", resp.StatusCode, err.Error()))
	}

	return string(ret), nil
}

func (c Client) SuccessRpc(class string, method string, args ...interface{}) (string, error) {
	ret, err := c.Send(class, method, args...)
	if err != nil {
		return "", err
	}

	type RetInfo struct {
		Code int `code`
		Data interface{} `data`
		Error string `error`
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
