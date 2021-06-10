package bind_empty

import (
	"fmt"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func Publish() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs2",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	var i = 0
	for {
		fmt.Printf("push: %d\n", i)

		err = ch.Publish(
			"logs2", // exchange
			"rk",    // routing key
			false,   // mandatory
			false,   // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(strconv.Itoa(i)),
			})
		PanicErr(err)
		i++
		time.Sleep(time.Second)
	}
}

func Consume() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs2",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	q, err := ch.QueueDeclare(
		"qd",
		false,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	ch.QueueBind(
		q.Name,
		"rk", // routing key
		"logs2",
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

func Consume2() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs2",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	q, err := ch.QueueDeclare(
		"qd2",
		false,
		false,
		false,
		false,
		nil,
	)
	PanicErr(err)

	ch.QueueBind(
		q.Name,
		"rk", // routing key
		"logs2",
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
