package main

// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

import (
	"github.com/garyburd/redigo/redis"
)

// Make a redis pool
var RedisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "pool:6379")
	},
}
