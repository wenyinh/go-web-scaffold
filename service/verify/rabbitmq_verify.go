package verify

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go-web-scaffold/mq"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	testQueueName = "test_queue"
	consumerCount = 5
	testMsgCount  = 10
)

func TestRabbitMQ() (err error) {
	// 使用全局发布通道
	ch := mq.PubChan
	if ch == nil {
		return fmt.Errorf("publish channel is nil")
	}
	// 声明队列，保证幂等
	_, err = ch.QueueDeclare(testQueueName, false, true, false, false, nil)
	if err != nil {
		zap.L().Error("QueueDeclare failed", zap.Error(err))
		return err
	}
	var msgWg sync.WaitGroup
	msgWg.Add(testMsgCount)
	// 启动多个消费者
	for i := 0; i < consumerCount; i++ {
		// 使用goroutine让每个消费者异步执行
		go func(workerId int) {
			if err = consume(testQueueName, workerId, &msgWg); err != nil {
				zap.L().Error("consumer failed", zap.Int("workerId", workerId), zap.Error(err))
			}
		}(i)
	}
	// 发布消息
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("msg_%d", i)
		if err = mq.PublishMessage(testQueueName, msg); err != nil {
			zap.L().Error("PublishMessage failed", zap.Error(err), zap.String("msg", msg))
			return err
		}
	}
	fmt.Println("Publish messages success")
	msgWg.Wait()
	fmt.Println("All test messages processed successfully, start web server...")
	return
}

func consume(queueName string, workerId int, msgWg *sync.WaitGroup) error {
	ch, err := mq.Conn.Channel() // 每调用一次，创建一个不同的channel，每个消费者有一个独立的channel
	if err != nil {
		return fmt.Errorf("create channel failed: %w", err)
	}
	defer ch.Close()
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume failed: %w", err)
	}
	for msg := range msgs {
		processMessage(msg, workerId, msgWg)
	}
	return errors.New("consumer channel exit unexpectedly")
}

func processMessage(msg amqp.Delivery, workerId int, msgWg *sync.WaitGroup) {
	fmt.Printf("[Worker %d] Received: %s\n", workerId, string(msg.Body))
	time.Sleep(1 * time.Second)
	if err := msg.Ack(false); err != nil {
		zap.L().Error("Ack failed", zap.Error(err))
	}
	msgWg.Done()
}
