package cache

import (
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

var cacheClient Client

func init() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	redisOptions, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err.Error())
	}

	cacheClient = NewRedisCache(redis.NewClient(redisOptions))
}

func TestRedisCacheListPlain(t *testing.T) {
	cacheKey := "list_plain"
	if err := cacheClient.Set(cacheKey, []string{"test"}, time.Minute); err != nil {
		t.Error(err)
	}

	listTest := []string{}
	if err := cacheClient.Get(cacheKey, &listTest); err != nil {
		t.Error(err)
	}

	if len(listTest) == 0 {
		t.Errorf("empty list")
	}
}

func TestRedisCacheListObject(t *testing.T) {
	cacheKey := "list_object"
	if err := cacheClient.Set(cacheKey, []map[string]string{{
		"foo": "bar",
	}}, time.Minute); err != nil {
		t.Error(err)
	}

	listTest := []map[string]string{}
	if err := cacheClient.Get(cacheKey, &listTest); err != nil {
		t.Error(err)
	}

	if len(listTest) == 0 {
		t.Errorf("empty list")
	}
}
