package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go-web-scaffold/dao/mysql"
	"go-web-scaffold/dao/redis"
	"go-web-scaffold/logger"
	"go-web-scaffold/logic/verify"
	"go-web-scaffold/mq"
	"go-web-scaffold/router"
	"go-web-scaffold/settings"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. 加载配置
	fmt.Println("Starting Web Scaffold...")
	if err := settings.Init(); err != nil {
		fmt.Println("Init Settings Error:", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Println("Init logger Error:", err)
		return
	}
	zap.L().Debug("zap init success")
	defer zap.L().Sync()
	// 3. 初始化MySQL连接
	if err := mysql.Init(); err != nil {
		fmt.Println("Init mysql Error:", err)
		return
	}
	// 4. 验证MySQL服务
	if err := verify.TestMySQL(); err != nil {
		fmt.Println("Test mysql Error:", err)
		return
	}
	defer mysql.CloseMySQL()
	// 5. 初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Println("Init redis Error:", err)
		return
	}
	// 6. 验证Redis服务
	if err := verify.TestRedis(); err != nil {
		fmt.Println("Test redis Error:", err)
		return
	}
	defer redis.CloseRedis()
	// 7. 初始化RabbitMQ连接
	if err := mq.InitRabbitMQ(); err != nil {
		fmt.Println("Init rabbitmq Error:", err)
		return
	}
	// 8. 验证RabbitMQ服务
	if err := verify.TestRabbitMQ(); err != nil {
		fmt.Println("Test rabbitmq Error:", err)
		return
	}
	defer mq.CloseRabbitMQ()
	// 9. 注册路由
	r := router.Setup()
	// 10. 关机设置
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
