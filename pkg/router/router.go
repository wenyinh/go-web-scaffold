package router

import (
	"github.com/gin-gonic/gin"
	"go-web-scaffold/api/controller"
	"go-web-scaffold/pkg/logger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册路由
	r.GET("/ping", controller.Ping)
	r.GET("/user", controller.GetByRange)               //查询参数
	r.POST("/user/create", controller.CreateUser)       //请求体参数
	r.PUT("/user/update", controller.UpdateUser)        //请求体参数
	r.DELETE("/user/delete/:id", controller.DeleteUser) // 路径参数
	r.POST("/user/save/:id", controller.SaveUser)       // 路径参数➕请求体参数
	return r
}
