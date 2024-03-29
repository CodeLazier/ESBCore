package brokers

import (
	"github.com/go-redis/redis"
	"common/task/drive"
	"common/task/message"
	"common/task/util/yjson"
	"common/task/yerrors"
	"time"
)

type RedisBroker struct {
	client   *drive.RedisClient
	host     string
	port     string
	password string
	db       int
	poolSize int
}

func NewRedisBroker(host string, port string, password string, db int, poolSize int) RedisBroker {
	return RedisBroker{
		host:     host,
		port:     port,
		password: password,
		db:       db,
		poolSize: poolSize,
	}
}

func (r *RedisBroker) Activate() {
	client := drive.NewRedisClient(r.host, r.port, r.password, r.db, r.poolSize)
	r.client = &client
}

func (r *RedisBroker) SetPoolSize(n int) {
	r.poolSize = n
}
func (r *RedisBroker) GetPoolSize() int {
	return r.poolSize
}

func (r *RedisBroker) Next(queryName string) (message.Message, error) {
	var msg message.Message
	values, err := r.client.BLPop(queryName, 2*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return msg, yerrors.ErrEmptyQuery{}
		}
		return msg, err
	}

	err = yjson.YJson.UnmarshalFromString(values[1], &msg)
	return msg, err
}

func (r *RedisBroker) Send(queryName string, msg message.Message) error {
	b, err := yjson.YJson.Marshal(msg)

	if err != nil {
		return err
	}
	err = r.client.RPush(queryName, b)
	return err
}
