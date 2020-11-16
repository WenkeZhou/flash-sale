package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// 连接信息
const MQURL = "amqp://guest:guest@localhost:5672/"

// rabbitMQ 连接体
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string // 队列名称
	Exchange  string // 交换机名称
	Key       string // bind key 名称
	MqUrl     string // 连接信息
}

// 创建RabbitMQ结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MqUrl: MQURL}
}

// 断开 channel 和 connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("[err]%s, [message]%s", err, message)
	}
}

// 创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rmq := NewRabbitMQ(queueName, "", "")
	var err error
	// 获取 connection
	rmq.conn, err = amqp.Dial(rmq.MqUrl)
	rmq.failOnErr(err, "创建连接错误")
	rmq.channel, err = rmq.conn.Channel()
	rmq.failOnErr(err, "创建连接错误")
	return rmq
}

func (r *RabbitMQ) PublishSimple(message string) {
	// 1.申请队列, 如果队列不存在会自动创建, 如果存在则跳过创建
	// 保证队列存在, 消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		// name
		r.QueueName,
		// durable, 是否持久化
		false,
		// autoDelete, 是否为自动化删除
		false,
		// exclusive, 是否具有排他性
		false,
		// noWait, 是否阻塞
		false,
		// args,
		nil,
	)

	if err != nil {
		fmt.Println(err)
	}

	// 2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true, 根据 exchange 类型和 routekey 规则,
		// 如果无法找到该队列，那么会把发送的消息返回给发送者。
		false,
		// 如果为ture时, 当 exchange 发送消息到队列后发现队列上没有绑定消费者,
		// 这会把消息返回给消费者。
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (r *RabbitMQ) ConsumeSimple() {
	// 1.申请队列, 如果队列不存在会自动创建, 如果存在则跳过创建
	// 保证队列存在, 消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否为自动化删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答, 默认为true; 为false 需要手动实现反馈回调
		true,
		// 是否具有排他性,
		false,
		// 如果设置为true, 表示不能将同一个conenction中发送的消息传递给这个connection中的消费者
		false,
		// 队列消费是否阻塞, false 表示阻塞。
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)

	// 启动协程处理消息
	go func() {
		for d := range msgs {
			// 实现我们要处理的逻辑函数
			log.Printf("Receive a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message. To exit to press CTRL+C")
	<-forever
}

func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rmq := NewRabbitMQ("", exchangeName, "")
	var err error
	// 获取 connection
	rmq.conn, err = amqp.Dial(rmq.MqUrl)
	rmq.failOnErr(err, "创建连接错误")
	rmq.channel, err = rmq.conn.Channel()
	rmq.failOnErr(err, "创建连接错误")
	return rmq
}

// 订阅生成模式
func (r *RabbitMQ) PublishPub(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare exchange.")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	r.failOnErr(err, "Failed to publish message.")
}

func (r *RabbitMQ) ReceiveSub() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare exchange")

	// 2.试探创建队列
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to decalre queue")

	// 绑定队列到 exchange中
	err = r.channel.QueueBind(
		q.Name,
		// 在pub/sub模式下, 这里的key必须为空
		"",
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to decalre queue")

	// 消费信息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
		}
	}()

	<-forever
}

// 路由模式
// 创建 RabbitMQ 实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	rmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取 connection
	rmq.conn, err = amqp.Dial(rmq.MqUrl)
	rmq.failOnErr(err, "创建连接错误")
	rmq.channel, err = rmq.conn.Channel()
	rmq.failOnErr(err, "创建连接错误")
	return rmq
}

// 路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare  exchange")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	r.failOnErr(err, "Failed to declare  exchange")
}

// 路由模式接受消息
func (r *RabbitMQ) ReceiveRouting() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare  exchange")

	// 2.创建队列, 这里队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", // 随机产生队列的名字
		false,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare Queue")

	// 3.绑定队列
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to bind queue")

	// 4.消费信息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}

// 话题模式
// 创建 RabbitMQ 实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	rmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取 connection
	rmq.conn, err = amqp.Dial(rmq.MqUrl)
	rmq.failOnErr(err, "创建连接错误")
	rmq.channel, err = rmq.conn.Channel()
	rmq.failOnErr(err, "创建连接错误")
	return rmq
}

// 话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare  exchange")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	r.failOnErr(err, "Failed to declare  exchange")
}

// 话题模式接受消息
func (r *RabbitMQ) ReceiveTopic() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare  exchange")

	// 2.创建队列, 这里队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", // 随机产生队列的名字
		false,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare Queue")

	// 3.绑定队列
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to bind queue")

	// 4.消费信息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}
