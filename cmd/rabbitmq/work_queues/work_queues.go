package work_queues

import (
	"bytes"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func Produce(body []byte) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	PanicErr(err)

	err = ch.Qos(1, 0, false)
	PanicErr(err)

	err = ch.Publish("",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body},
	)
	PanicErr(err)
	fmt.Printf("send: %s\n", body)
}

func Consume() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	PanicErr(err)

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

	go func() {
		for msg := range msgs {
			fmt.Printf("receive: %s\n", msg.Body)

			dotCount := bytes.Count(msg.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)

			fmt.Printf("done\n")
		}
	}()

	<-stop
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
