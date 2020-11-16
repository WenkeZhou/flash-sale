package main

import (
	"fmt"
	"github.com/WenkeZhou/flash-sale/pkg/rabbitmq"
	"strconv"
	"time"
)

func main() {
	rmq1 := rabbitmq.NewRabbitMQTopic("exchange_topic", "keya.keyb.keyc")
	rmq2 := rabbitmq.NewRabbitMQTopic("exchange_topic", "keya.keyb.key2")
	for i := 0; i < 10; i++ {
		rmq1.PublishTopic("[keya.keyb.keyc]:" + strconv.Itoa(i))
		rmq2.PublishTopic("[keya.keyb.key2]" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
