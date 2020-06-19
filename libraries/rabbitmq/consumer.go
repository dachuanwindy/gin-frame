package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util"
)

type TestConsumer struct {}

//消费mq消息
func (self *TestConsumer) Do(d amqp.Delivery, header *log.LogFormat) error {
	fmt.Println(string(d.Body))
	err := d.Ack(false)
	util.Must(err)
	return err
}
