package utils

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func InitConnect() *amqp091.Connection {
	conn, err := amqp091.Dial("amqp://lian:020109@59.110.54.159:5672/")
	if err != nil {
		log.Print("mq connect err")
		log.Print(err)
		return conn
	}
	return conn
}
