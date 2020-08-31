package errors

import (
	"github.com/gopusher/gateway/pkg/web/response"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.LogError(c, err)
				response.ErrorJson(c, response.CodeSystemError, "")
			}
		}()
		c.Next()
	}

	//return gin.Recovery()
}
