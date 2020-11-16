package main

import (
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
)

func main() {
	rmq := rabbitmq.NewRabbitMQPubSub("exchange_test")
	rmq.ReceiveSub()
}
