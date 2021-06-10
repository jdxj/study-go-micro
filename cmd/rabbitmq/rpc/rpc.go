package rpc

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

func Publish() {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue2",
		false,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	err = ch.Qos(1, 0, false)
	PanicErr(err)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	stop := make(chan int)

	go func() {
		for msg := range msgs {
			n, _ := strconv.Atoi(string(msg.Body))
			resp := fib(n)

			err = ch.Publish(
				"",
				msg.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: msg.CorrelationId,
					Body:          []byte(strconv.Itoa(resp)),
				})
			PanicErr(err)

			msg.Ack(false)
		}
	}()

	<-stop

}

func Consume() {
	n := os.Args[2]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	PanicErr(err)

	fmt.Printf("consumer queue name: %s\n", q.Name)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	id := strconv.Itoa(rand.Int())

	err = ch.Publish(
		"",
		"rpc_queue2",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: id,
			ReplyTo:       q.Name,
			Body:          []byte(n),
		},
	)
	PanicErr(err)

	for msg := range msgs {
		resp, _ := strconv.Atoi(string(msg.Body))
		fmt.Printf("%d\n", resp)
		break
	}
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}

	return fib(n-1) + fib(n-2)
}
