package router

import (
	"github.com/gin-gonic/gin"
	"go-web-scaffold/dao/redis"
	"go-web-scaffold/logger"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册路由
	r.GET("/ping", func(c *gin.Context) {
		redis.RDB.Set("Hello", "Golang", 0)
		c.String(http.StatusOK, "PONG!")
	})
	return r
}
