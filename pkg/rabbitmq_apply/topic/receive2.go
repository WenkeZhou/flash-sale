package main

import (
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
)

func main() {
	rmq := rabbitmq.NewRabbitMQTopic("exchange_topic", "keya.*.key2")
	rmq.ReceiveTopic()
}
