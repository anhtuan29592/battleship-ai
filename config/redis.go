package config

import (
	"github.com/go-redis/redis"
)

var DefaultRedis *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})
