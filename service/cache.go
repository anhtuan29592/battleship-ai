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

	err = c.RedisCli.Set(key, string(arr), time.Duration(5 * time.Minute)).Err()
	if err != nil {
		log.Print("put cache error, retry...")
		err = c.RedisCli.Set(key, string(arr), time.Duration(5 * time.Minute)).Err()
	}

	return err
}

func (c *CacheService) Get(key string, out interface{}) error {

	jsonVal, err := c.RedisCli.Get(key).Result()
	if err != nil {
		log.Print("get cache error, retry...")
		count := 0
		for {
			log.Printf("retry %d time", count)
			if count > 10 {
				break
			}
			jsonVal, err = c.RedisCli.Get(key).Result()
			if err == nil {
				break
			}
			count++
			time.Sleep(1 * time.Second)
		}
	}

	if err != nil {
		return err
	}
	json.Unmarshal([]byte(jsonVal), &out)
	return nil
}
