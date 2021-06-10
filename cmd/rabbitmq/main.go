package main

import (
	"os"
	"study_go_micro/cmd/rabbitmq/bind_empty"
)

func main() {
	switch os.Args[1] {
	case "1":
		bind_empty.Publish()
	case "2":
		bind_empty.Consume2()
	}
}
