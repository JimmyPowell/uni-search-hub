package middleware

import (
	"github.com/gin-gonic/gin"
)

// TODO 后续需要改进
func abortWithMessage(c *gin.Context, statusCode int, message string, code ...string) {
	codeStr := ""
	if len(code) > 0 {
		codeStr = code[0]
	}
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"message": message,
			"type":    "new_api_error",
			"code":    codeStr,
		},
	})
	c.Abort()
}
