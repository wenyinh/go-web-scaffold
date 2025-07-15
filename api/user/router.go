package user

import (
	"github.com/gin-gonic/gin"
	"go-web-scaffold/pkg/logger"
)

func RegisterRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册路由
	r.GET("/ping", Ping)
	r.GET("/user", GetByRange)               //查询参数
	r.POST("/user/create", CreateUser)       //请求体参数
	r.PUT("/user/update", UpdateUser)        //请求体参数
	r.DELETE("/user/delete/:id", DeleteUser) // 路径参数
	r.POST("/user/save/:id", SaveUser)       // 路径参数➕请求体参数
	return r
}
