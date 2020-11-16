package main

import (
	"fmt"
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rmq := rabbitmq.NewRabbitMQPubSub("exchange_test")
	for i := 1; i <= 100; i++ {
		rmq.PublishPub("this is pub message: " + strconv.Itoa(i))
		fmt.Printf("this is pub message: " + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
	}
}
