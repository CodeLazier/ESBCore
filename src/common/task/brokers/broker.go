package brokers

import (
	"common/task/message"
)

type BrokerInterface interface {
	Next(queryName string) (message.Message, error)
	Send(queryName string, msg message.Message) error
	// 调用Activate后才真正建立连接
	Activate()
	SetPoolSize(int)
	GetPoolSize()int
}
