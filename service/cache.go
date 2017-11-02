package service

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
	"log"
)

type CacheService struct {
	RedisCli *redis.Client `inject:""`
}

func (c *CacheService) Put(key string, value interface{}, duration int) error {
	arr, err := json.Marshal(value)
	if err != nil {
		log.Print("json marshal error, retry...")
		arr, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	err = c.RedisCli.Set(key, string(arr), time.Duration(duration)).Err()
	if err != nil {
		log.Print("put cache error, retry...")
		err = c.RedisCli.Set(key, string(arr), time.Duration(duration)).Err()
	}

	return err
}

func (c *CacheService) Get(key string, out interface{}) error {
	jsonVal, err := c.RedisCli.Get(key).Result()
	if err != nil {
		log.Print("get cache error, retry...")
		jsonVal, err = c.RedisCli.Get(key).Result()
		if err != nil {
			return err
		}
	}
	json.Unmarshal([]byte(jsonVal), &out)
	return nil
}
