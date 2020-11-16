package main

import (
	"fmt"
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
)

func main() {
	rmq := rabbitmq.NewRabbitMQSimple("simple")
	rmq.ConsumeSimple()
	fmt.Println("发送成功")
}
