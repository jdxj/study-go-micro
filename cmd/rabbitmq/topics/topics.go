package topics

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func Publish() {
	// program name, routing key, body
	routingKey := os.Args[1]
	body := []byte(os.Args[2])

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	err = ch.Publish(
		"logs_topic", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

func Consume() {
	// program name, routing key
	routingKey := os.Args[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	PanicErr(err)

	ch.QueueBind(
		q.Name,
		routingKey, // routing key
		"logs_topic",
		false,
		nil,
	)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	stop := make(chan int)

	for msg := range msgs {
		fmt.Printf("%s\n", msg.Body)
	}

	<-stop
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
