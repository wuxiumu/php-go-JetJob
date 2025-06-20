package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// TokenAuthMiddleware 用于保护API，校验 Bearer token
func Cors(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"message": "token格式错误"})
			return
		}
		userToken := strings.TrimSpace(h[7:])
		if userToken != token {
			c.AbortWithStatusJSON(401, gin.H{"message": "token无效"})
			return
		}
		c.Next()
	}
}
