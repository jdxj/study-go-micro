package pub_sub

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func Publish() {

	body := ""
	if len(os.Args) > 1 {
		body = os.Args[1]
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
}

func Consume() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
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
		"",
		"logs",
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
