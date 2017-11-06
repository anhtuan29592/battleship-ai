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

func (c *CacheService) Put(key string, value interface{}) error {
	arr, err := json.Marshal(value)
	if err != nil {
		log.Print("json marshal error, retry...")
		arr, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	err = c.RedisCli.Set(key, string(arr), time.Duration(1800000000000)).Err()
	if err != nil {
		log.Print("put cache error, retry...")
		err = c.RedisCli.Set(key, string(arr), time.Duration(1800000000000)).Err()
	}

	return err
}

func (c *CacheService) Get(key string, out interface{}) error {

	count := 5
	jsonVal, err := c.RedisCli.Get(key).Result()
	if err != nil {
		log.Print("get cache error, retry...")
		for {
			jsonVal, err = c.RedisCli.Get(key).Result()
			if err == nil {
				break
			}

			if count < 0 {
				break
			}

			count--
		}
	}

	if err != nil {
		return err
	}
	json.Unmarshal([]byte(jsonVal), &out)
	return nil
}
