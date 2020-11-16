package main

import (
	"fmt"
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
)

func main() {
	rmq := rabbitmq.NewRabbitMQSimple("simple")
	rmq.PublishSimple("hello this is at test!")
	fmt.Println("发送成功")
}
