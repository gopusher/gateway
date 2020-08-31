package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Token  string `mapstructure:"token" validate:"required"`
	Secret string `mapstructure:"secret"`
}

func NewRobot(config *Config) *Robot {
	return &Robot{
		token:  config.Token,
		secret: config.Secret,
	}
}

func sign(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

type Robot struct {
	token, secret string
}

func (robot *Robot) SendMessage(msg interface{}) error {
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(msg)
	if err != nil {
		return fmt.Errorf("msg json failed, msg: %v, err: %v", msg, err.Error())
	}

	value := url.Values{}
	value.Set("access_token", robot.token)
	if robot.secret != "" {
		t := time.Now().UnixNano() / 1e6
		value.Set("timestamp", fmt.Sprintf("%d", t))
		value.Set("sign", sign(t, robot.secret))
	}

	request, err := http.NewRequest(http.MethodPost, "https://oapi.dingtalk.com/robot/send", body)
	if err != nil {
		return fmt.Errorf("error request: %v", err.Error())
	}
	request.URL.RawQuery = value.Encode()
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	res, err := (&http.Client{}).Do(request)
	if err != nil {
		return fmt.Errorf("send dingTalk message failed, error: %v", err.Error())
	}
	defer func() { _ = res.Body.Close() }()
	result, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return fmt.Errorf("send dingTalk message failed, %s", httpError(request, res, result, "http code is not 200"))
	}
	if err != nil {
		return fmt.Errorf("send dingTalk message failed, %s", httpError(request, res, result, err.Error()))
	}

	type response struct {
		ErrCode int `json:"errcode"`
	}
	var ret response

	if err := json.Unmarshal(result, &ret); err != nil {
		return fmt.Errorf("send dingTalk message failed, %s", httpError(request, res, result, err.Error()))
	}

	if ret.ErrCode != 0 {
		return fmt.Errorf("send dingTalk message failed, %s", httpError(request, res, result, "errcode is not 0"))
	}

	return nil
}

func httpError(request *http.Request, response *http.Response, body []byte, error string) string {
	return fmt.Sprintf(
		"http request failure, error: %s, status code: %d, %s %s, body:\n%s",
		error,
		response.StatusCode,
		request.Method,
		request.URL.String(),
		string(body),
	)
}

func (robot *Robot) SendTextMessage(content string, atMobiles []string, isAtAll bool) error {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
		"at": map[string]interface{}{
			"atMobiles": atMobiles,
			"isAtAll":   isAtAll,
		},
	}

	return robot.SendMessage(msg)
}

func (robot *Robot) SendMarkdownMessage(title string, text string, atMobiles []string, isAtAll bool) error {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
		"at": map[string]interface{}{
			"atMobiles": atMobiles,
			"isAtAll":   isAtAll,
		},
	}

	return robot.SendMessage(msg)
}

func (robot *Robot) SendLinkMessage(title string, text string, messageUrl string, picUrl string) error {
	msg := map[string]interface{}{
		"msgtype": "link",
		"link": map[string]string{
			"title":      title,
			"text":       text,
			"messageUrl": messageUrl,
			"picUrl":     picUrl,
		},
	}

	return robot.SendMessage(msg)
}
