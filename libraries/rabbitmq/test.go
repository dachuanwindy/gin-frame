package rabbitmq

import (
	"gin-frame/libraries/util"
)

func consumer() {
	connUri := "amqp://why:why@localhost:5672/why"
	queueName := "why_queue"

	testConsumer := &TestConsumer{}
	c, err := NewConsumer(connUri, queueName, "", testConsumer)
	util.Must(err)
	defer c.Shutdown()
}

func producer() {
	connUri := "amqp://why:why@localhost:5672/why"
	queueName := "why_queue"
	exchangeName := "why_exchange"
	routeName := "why_route"

	err := NewProducerCmd("hello world!", connUri, exchangeName, "direct", queueName, routeName, "")
	util.Must(err)
}

func main() {
	producer()
	consumer()
}
