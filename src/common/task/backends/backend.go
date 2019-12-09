package backends

import "common/task/message"

type BackendInterface interface {
	SetResult(result message.Result, exTime int) error
	GetResult(key string) (message.Result, error)
	// 调用Activate后才真正建立连接
	Activate()
	SetPoolSize(int)
	GetPoolSize() int
}
