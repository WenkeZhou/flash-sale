package main

import (
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
)

func main() {
	rmq := rabbitmq.NewRabbitMQSimple("simple_work_queue")
	rmq.ConsumeSimple()
}
