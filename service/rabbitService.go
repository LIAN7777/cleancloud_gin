package service

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func StartConsumer(name string, handler func(msg []byte)) {
	conn, err := amqp091.Dial("amqp://lian:020109@59.110.54.159:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()
	wg.Wait()
}

func PrintComment(msg []byte) {
	fmt.Print(string(msg))
}
