package rabit

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Init() error {
	// 建立与 RabbitMQ 服务器的连接
	conn, err := amqp.Dial("amqp://admin:admin123@124.221.8.185:5672")
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	//// 声明队列
	//queueName := "my_queue"
	//_, err = ch.QueueDeclare(
	//	queueName, // 队列名称
	//	false,     // 是否持久化
	//	false,     // 是否自动删除
	//	false,     // 是否独占模式
	//	false,     // 是否阻塞等待
	//	nil,       // 额外参数
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 发布消息
	//message := "Hello, RabbitMQ!"
	//err = ch.Publish(
	//	"",        // 交换机名称
	//	queueName, // 路由键
	//	false,     // 是否强制
	//	false,     // 是否立即
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte(message),
	//	},
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 消费消息
	//msgs, err := ch.Consume(
	//	queueName, // 队列名称
	//	"",        // 消费者名称
	//	true,      // 是否自动应答
	//	false,     // 是否独占模式
	//	false,     // 是否阻塞等待
	//	false,     // 是否为非阻塞模式
	//	nil,       // 额外参数
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 处理接收到的消息
	//for msg := range msgs {
	//	fmt.Println(string(msg.Body))
	//}

	return nil
}
