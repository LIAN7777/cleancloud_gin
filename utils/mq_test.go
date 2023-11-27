package utils

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"strconv"
	"testing"
)

func TestInitConnect(t *testing.T) {
	channel, err := InitConnect().Channel()
	if err != nil {
		fmt.Print("init fail")
		return
	}
	for i := 0; i < 10; i++ {
		err = channel.Publish(
			"",           // exchange
			"test_queue", // routing key
			false,        // mandatory
			false,        // immediate
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte("this is a message" + strconv.Itoa(i)),
			},
		)
		if err != nil {
			fmt.Print("message publish fail")
			return
			// 处理发布消息错误
		}
	}
}

func TestConsume(t *testing.T) {
	channel, err := InitConnect().Channel()
	if err != nil {
		fmt.Print("init fail")
		return
	}
	msgs, err := channel.Consume(
		"test_queue", // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		// 处理消费者创建错误
	}

	go func() {
		for msg := range msgs {
			// 处理接收到的消息
			message := string(msg.Body)
			fmt.Println("Received message:", message)
		}
	}()
}
func TestConsumeOne(t *testing.T) {
	channel, err := InitConnect().Channel()
	if err != nil {
		fmt.Print("init fail")
		return
	}
	msg, ok, err := channel.Get("test_queue", true)
	if err != nil {
		//消息获取错误
	}
	if ok {
		message := string(msg.Body)
		fmt.Println("Received message:", message)
	}
}
