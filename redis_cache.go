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

func (c *redisCache) Nearby(k string, lon, lat, radius float64) ([]Location, error) {
	result := []Location{}
	locations, err := c.client.GeoRadius(k, lon, lat, &redis.GeoRadiusQuery{
		Radius:   radius,
		Unit:     "m",
		WithDist: true,
		Sort:     "ASC",
	}).
		Result()
	if err != nil {
		return nil, err
	}

	for _, location := range locations {
		result = append(result, Location{
			Name:      location.Name,
			Longitude: location.Longitude,
			Latitude:  location.Latitude,
			Distance:  location.Dist,
		})
	}

	return result, nil
}

func (c *redisCache) GeoAdd(k string, locations ...Location) error {
	redisLocations := []*redis.GeoLocation{}
	for _, location := range locations {
		redisLocations = append(redisLocations, &redis.GeoLocation{
			Name:      location.Name,
			Longitude: location.Longitude,
			Latitude:  location.Latitude,
		})
	}

	return c.client.GeoAdd(k, redisLocations...).Err()
}

func (c *redisCache) HSet(k string, f string, v interface{}, expire time.Duration) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}

	success := c.client.HSet(k, f, buf).Err()
	c.client.Expire(k, expire)

	return success
}

func (c *redisCache) HGet(k string, f string, v interface{}) error {

	rv, err := c.client.HGet(k, f).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rv), v); err != nil {
		return err
	}

	return nil
}

func NewRedisCache(redisClient *redis.Client) Client {
	return &redisCache{
		client: redisClient,
	}
}
