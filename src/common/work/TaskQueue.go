package work

import (
	"common/task"
	"common/task/backends"
	"common/task/brokers"
	"common/task/server"

)

func NewTaskBroker(addr string, port string, password string, db int, pool int) *brokers.RedisBroker {
	b := task.Broker.NewRedisBroker(addr, port, password, db, pool)
	return &b
}

func NewTaskBackend(addr string, port string, password string, db int, pool int) *backends.RedisBackend {
	b := task.Backend.NewRedisBackend(addr, port, password, db, pool)
	return &b
}

func InitTaskServer(broker *brokers.RedisBroker, backend *backends.RedisBackend, statusExpires int, resultExpires int) *server.Server {

	s := task.Server.NewServer(
		task.Config.Broker(broker),
		task.Config.Backend(backend),
		task.Config.StatusExpires(statusExpires),
		task.Config.ResultExpires(resultExpires),
	)

	return &s
}
