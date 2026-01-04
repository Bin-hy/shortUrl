package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 捕获所有 panic，防止服务崩溃
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印错误堆栈信息
				log.Printf("panic recovered: %v\n%s", err, debug.Stack())

				// 返回 500 错误给客户端
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
