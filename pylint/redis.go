package pylint

// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// Make a redis pool
var RedisPool *redis.Pool

func ConnectRedis(cfg Config) redis.Conn {
	RedisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			connInfo := fmt.Sprintf(
				"%s:%s",
				cfg.Redis.Host,
				cfg.Redis.Port,
			)
			return redis.Dial("tcp", connInfo)
		},
	}
	conn := RedisPool.Get()
	return conn
}
