package auth

import (
	"github.com/gopusher/gateway/pkg/web/response"

	"github.com/gin-gonic/gin"
)

func Check(authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.GetHeader("X-API-TOKEN"); token != authToken {
			response.ErrorJson(c, response.CodeUnauthorized, "token error")
			c.Abort()
			return
		}
		c.Next()
	}
}
