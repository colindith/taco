package broker

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }

func Send() {
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqUser := os.Getenv("RABBITMQ_DEFAULT_USER")
	rabbitmqPass := os.Getenv("RABBITMQ_DEFAULT_PASS")
	rabbitmqVhost := os.Getenv("RABBITMQ_DEFAULT_VHOST")
	rabbitmqDial := fmt.Sprintf("amqp://%s:%s@%s:5672/%s", rabbitmqUser, rabbitmqPass, rabbitmqHost, rabbitmqVhost)
	fmt.Println("----------- ", rabbitmqDial)
	conn, err := amqp.Dial(rabbitmqDial)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
