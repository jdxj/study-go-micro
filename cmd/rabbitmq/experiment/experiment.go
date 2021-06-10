package experiment

import (
	"fmt"
	"strconv"
	"study_go_micro/cmd/rabbitmq/util"
	"time"

	"github.com/streadway/amqp"
)

const (
	exchangeName = "experiment"
)

func Producer() {
	conn := util.NewRabbitmqConn()

	mqChan, err := conn.Channel()
	util.PanicErr(err)
	defer mqChan.Close()

	err = mqChan.ExchangeDeclare(
		exchangeName,
		"topic",
		false,
		false,
		false,
		false,
		nil,
	)
	util.PanicErr(err)

	i := 0
	for {
		err = mqChan.Publish(
			exchangeName,
			"abc", // 推送的时候可以附带 routing key
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(strconv.Itoa(i)),
			},
		)
		util.PanicErr(err)

		fmt.Printf("push access: %d\n", i)
		time.Sleep(5 * time.Second)
		i++
	}
}

func Consumer() {
	conn := util.NewRabbitmqConn()

	mqChan, err := conn.Channel()
	util.PanicErr(err)
	defer mqChan.Close()

	q, err := mqChan.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	util.PanicErr(err)

	err = mqChan.QueueBind(
		q.Name,
		"abc",
		exchangeName,
		false,
		nil,
	)
	util.PanicErr(err)

	msgs1, err := mqChan.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	util.PanicErr(err)

	go func() {
		for msg := range msgs1 {
			fmt.Printf("msg1: %s\n", msg.Body)
		}
	}()

	msgs2, err := mqChan.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	util.PanicErr(err)

	go func() {
		for msg := range msgs2 {
			fmt.Printf("msg2: %s\n", msg.Body)
		}
	}()

	time.Sleep(time.Hour)
}
