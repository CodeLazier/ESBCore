package work

import (
	. "common"
	"github.com/gojuukaze/YTask/v2"
	"github.com/gojuukaze/YTask/v2/backends"
	"github.com/gojuukaze/YTask/v2/brokers"
	"github.com/gojuukaze/YTask/v2/server"
)

func NewTaskBroker(addr string, port string, password string, db int, pool int) *brokers.RedisBroker {
	b := ytask.Broker.NewRedisBroker(addr, port, password, db, pool)
	return &b
}

func NewTaskBackend(addr string, port string, password string, db int, pool int) *backends.RedisBackend {
	b := ytask.Backend.NewRedisBackend(addr, port, password, db, pool)
	return &b
}

func InitTaskServer(broker *brokers.RedisBroker, backend *backends.RedisBackend, statusExpires int, resultExpires int) *server.Server {

	s := ytask.Server.NewServer(
		ytask.Config.Broker(broker),
		ytask.Config.Backend(backend),
		ytask.Config.Debug(true),
		ytask.Config.StatusExpires(statusExpires),
		ytask.Config.ResultExpires(resultExpires),
	)
	s.Add(ESBCaller, "MainEnter", MainEnter)


	return &s
}
