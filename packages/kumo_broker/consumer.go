package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)
// The listener daemon process for Kumo message
type Consumer struct {
  
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqUser := os.Getenv("RABBITMQ_DEFAULT_USER")
	rabbitmqPass := os.Getenv("RABBITMQ_DEFAULT_PASS")
	rabbitmqVhost := os.Getenv("RABBITMQ_DEFAULT_VHOST")
	rabbitmqDial := fmt.Sprintf("amqp://%s:%s@%s:5672/%s", rabbitmqUser, rabbitmqPass, rabbitmqHost, rabbitmqVhost)
	conn, err := amqp.Dial(rabbitmqDial)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"kumo_broadcast",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"kumo_broadcast", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

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

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			log.Printf(" [%s] receive the msg: %s", msg.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
