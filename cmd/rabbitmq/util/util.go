package util

import "github.com/streadway/amqp"

func NewRabbitmqConn() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	return conn
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
