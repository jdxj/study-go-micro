package hello_world

import (
	"fmt"

	"github.com/streadway/amqp"
)

func Produce() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	PanicErr(err)

	body := "hello rabbitmq"
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(body)})
	PanicErr(err)
}

func Consume() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	PanicErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	PanicErr(err)
	defer ch.Close()

	// 保证 queue 存在
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	PanicErr(err)

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	PanicErr(err)

	stop := make(chan int)
	go func() {
		for msg := range messages {
			fmt.Printf("%s\n", msg.Body)
		}
		close(stop)
	}()

	<-stop
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
