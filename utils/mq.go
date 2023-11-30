package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

var Channel *amqp091.Channel

func InitConnect() *amqp091.Connection {
	conn, err := amqp091.Dial("amqp://lian:020109@59.110.54.159:5672/")
	if err != nil {
		log.Print("mq connect err")
		log.Print(err)
		return conn
	}
	return conn
}

func InitRabbit() error {
	channel, err := InitConnect().Channel()
	if err != nil {
		fmt.Print("init fail")
		return err
	}
	Channel = channel
	return nil
}

func Publish(exchange string, queueOrRoute string, data interface{}) error {
	messageData, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = Channel.PublishWithContext(
		ctx,
		exchange,
		queueOrRoute,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        messageData,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
