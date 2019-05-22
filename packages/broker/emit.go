package broker

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func BroadcastMessage(msg string) {
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
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// body := bodyFrom(os.Args)
	// body := "broadcast msg ~~~ ~~~"
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", msg)
}

// func bodyFrom(args []string) string {
// 	var s string
// 	if (len(args) < 2) || os.Args[1] == "" {
// 		s = "hello"
// 	} else {
// 		s = strings.Join(args[1:], " ")
// 	}
// 	return s
// }
