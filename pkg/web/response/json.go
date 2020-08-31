package response

import (
	"github.com/gopusher/gateway/pkg/log"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func SuccessJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Code: CodeSuccess,
		Data: data,
		Msg:  "",
	})
}

func LogError(c *gin.Context, error interface{}) {
	request, _ := httputil.DumpRequest(c.Request, false)
	log.Error(
		"[Recovery] panic recovered",
		zap.Any("error", error),
		zap.String("request", string(request)),
	)
}

func ErrorJson(c *gin.Context, code int, msg string) {
	ErrorJsonWithStatusCode(c, http.StatusOK, code, msg)
}

func ErrorJsonWithStatusCode(c *gin.Context, statusCode, code int, msg string) {
	c.JSON(statusCode, &Response{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}
