package service

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type CacheService struct {
	RedisCli *redis.Client `inject:""`
}

func (c *CacheService) Put(key string, value interface{}, duration int) error {
	arr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.RedisCli.Set(key, string(arr), time.Duration(duration)).Err()
}

func (c *CacheService) Get(key string, out interface{}) error {
	jsonVal, err := c.RedisCli.Get(key).Result()
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(jsonVal), &out)
	return nil
}
