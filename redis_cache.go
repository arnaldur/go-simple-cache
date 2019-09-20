package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

type redisCache struct {
	client *redis.Client
}

func (c *redisCache) Delete(k string) error {
	return c.client.Del(k).Err()
}

func (c *redisCache) Flush() error {
	return c.client.FlushDB().Err()
}

func (c *redisCache) Get(k string, v interface{}) error {
	rv, err := c.client.Get(k).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rv), v); err != nil {
		return err
	}

	return nil
}

func (c *redisCache) Set(k string, v interface{}, expire time.Duration) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.client.Set(k, string(buf), expire).Err()
}

func NewRedisCache(redisClient *redis.Client) Client {
	return &redisCache{
		client: redisClient,
	}
}
