package notification

import (
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
		"bytes"
	"io/ioutil"
	"reflect"
	"github.com/gopusher/gateway/log"
)

type Client struct {
	url string
	userAgent string
}

func NewRpc(url string, userAgent string) *Client {
	log.Info("Notification url: %s, userAgent: %s", url, userAgent)

	return &Client {
		url: url,
		userAgent: userAgent,
	}
}

func (c Client) post(body []byte) (string, error) {
	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("notify failed, error: " + err.Error())
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("notify failed, error: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("notify failed, http code: %d, read body failed，error: %s", resp.StatusCode, err.Error()))
	}

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("notify failed, http code: %d, read body failed，error: %s", 200, err.Error()))
	}

	return string(ret), nil
}

func (c Client) Call(method string, args ...interface{}) (string, error) {
	type RpcBody struct {
		MethodName 	string	`json:"method"`
		Args		[]interface{}	`json:"args"`
	}
	body, err := json.Marshal(&RpcBody{
		MethodName: method,
		Args: args,
	})
	if err != nil {
		return "", errors.New("notify failed, error: " + err.Error())
	}

	// post data and get result.
	ret, err := c.post(body)

	if err != nil {
		return "", err
	}

	//parse result
	type RetInfo struct {
		Code int `json"code"`
		Data interface{} `json"data"`
		Error string `json"error"`
	}

	var retInfo RetInfo
	if err := json.Unmarshal([]byte(ret), &retInfo); err != nil {
		return "", errors.New(fmt.Sprintf("notify failed, parse body failed，type: %v, value: %s", reflect.TypeOf(ret), ret))
	}

	if retInfo.Code != 0 {
		return "", errors.New(retInfo.Error)
	}

	return ret, nil
}
