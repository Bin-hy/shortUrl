package interfaces

import (
	"github.com/Bin-hy/shortUrl/internal/interfaces/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *Handler) *gin.Engine {
	// 使用 New() 而不是 Default()，这样我们可以自定义中间件
	router := gin.New()
	
	// 添加日志中间件
	router.Use(gin.Logger())
	
	// 添加自定义的 Recovery 中间件
	router.Use(middleware.RecoveryMiddleware())

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", Ping)
		v1.POST("/shorten", h.ShortenV1)
		v1.GET("/:shortUrl", h.RedirectV1)
	}
	v2 := router.Group("/v2")
	{
		v2.GET("/ping", Ping)
		v2.POST("/shorten", h.ShortenV2)
		v2.GET("/:shortUrl", h.RedirectV2)
	}
	v3 := router.Group("/v3")
	{
		v3.GET("/ping", Ping)
		v3.POST("/shorten", h.ShortenV3)
		v3.GET("/:shortUrl", h.RedirectV3)
	}
	return router
}

func Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
