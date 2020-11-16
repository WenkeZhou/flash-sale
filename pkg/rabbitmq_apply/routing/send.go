package main

import (
	"fmt"
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rmq1 := rabbitmq.NewRabbitMQRouting("excahng_routing", "key_a")
	rmq2 := rabbitmq.NewRabbitMQRouting("excahng_routing", "key_b")
	for i := 0; i < 10; i++ {
		rmq1.PublishRouting("key_a:" + strconv.Itoa(i))
		rmq2.PublishRouting("key_b:" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
