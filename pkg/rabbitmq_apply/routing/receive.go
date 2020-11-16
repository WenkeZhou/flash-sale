package main

import (
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
)

func main() {
	rmq := rabbitmq.NewRabbitMQRouting("excahng_routing", "key_a")
	rmq.ReceiveRouting()
}
