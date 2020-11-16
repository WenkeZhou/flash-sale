package main

import (
	"fmt"
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rmq := rabbitmq.NewRabbitMQSimple("simple_work_queue")
	for i := 1; i <= 100; i++ {
		rmq.PublishSimple("message:" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
