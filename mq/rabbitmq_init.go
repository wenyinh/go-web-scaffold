package mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Conn *amqp.Connection

// PubChan 全局消息发布channel
var PubChan *amqp.Channel

func InitRabbitMQ() (err error) {
	zap.L().Info("init rabbitmq...")
	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetInt("rabbitmq.port"),
		viper.GetString("rabbitmq.vhost"),
	)
	Conn, err = amqp.Dial(url)
	if err != nil {
		zap.L().Error("rabbitmq init error", zap.Error(err))
		return
	}
	PubChan, err = Conn.Channel()
	if err != nil {
		zap.L().Error("rabbitmq channel creation error", zap.Error(err))
		return
	}
	fmt.Println("PONG! RabbitMQ Connect Success!")
	return
}

func CloseRabbitMQ() {
	if Conn != nil {
		_ = Conn.Close()
	}
	if PubChan != nil {
		_ = PubChan.Close()
	}
}

func PublishMessage(queueName string, msg string) (err error) {
	// 声明队列
	_, err = PubChan.QueueDeclare(queueName, false, true, false, false, nil)
	if err != nil {
		zap.L().Error("QueueDeclare failed: ", zap.Error(err))
		return
	}
	// 发布消息
	err = PubChan.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		zap.L().Error("Publish failed: ", zap.Error(err))
		return
	}
	fmt.Printf("Publish Success: %s\n", string(msg))
	return
}
