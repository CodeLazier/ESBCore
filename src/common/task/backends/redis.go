package backends

import (
	"github.com/go-redis/redis"
	"common/task/drive"
	"common/task/message"
	"common/task/util/yjson"
	"common/task/yerrors"
	"time"
)

type RedisBackend struct {
	client   *drive.RedisClient
	host     string
	port     string
	password string
	db       int
	poolSize int
}

func NewRedisBackend(host string, port string, password string, db int, poolSize int) RedisBackend {
	return RedisBackend{
		host:     host,
		port:     port,
		password: password,
		db:       db,
		poolSize: poolSize,
	}
}

func (r *RedisBackend) Activate() {
	client := drive.NewRedisClient(r.host, r.port, r.password, r.db, r.poolSize)
	r.client = &client
}

func (r *RedisBackend) SetPoolSize(n int) {
	r.poolSize = n
}
func (r *RedisBackend) GetPoolSize() int {
	return r.poolSize
}
func (r *RedisBackend) SetResult(result message.Result, exTime int) error {
	b, err := yjson.YJson.Marshal(result)

	if err != nil {
		return err
	}
	err = r.client.Set(result.GetBackendKey(), b, time.Duration(exTime)*time.Second)
	return err
}
func (r *RedisBackend) GetResult(key string) (message.Result, error) {
	var result message.Result

	b, err := r.client.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return result, yerrors.ErrNilResult{}
		}
		return result, err
	}

	err = yjson.YJson.Unmarshal(b, &result)
	return result, err
}
