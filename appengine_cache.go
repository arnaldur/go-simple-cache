package cache

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

type appengineCache struct {
	prefix string
	codec  memcache.Codec
	ctx    context.Context
}

func (c *appengineCache) key(k string) string { return c.prefix + "/" + k }

func (c *appengineCache) Delete(k string) error { return memcache.Delete(c.ctx, c.key(k)) }

func (c *appengineCache) Flush() error { return memcache.Flush(c.ctx) }

func (c *appengineCache) Get(k string, v interface{}) error {
	_, err := c.codec.Get(c.ctx, c.key(k), v)
	if err == memcache.ErrCacheMiss {
		return ErrCacheMiss
	}

	return err
}

func (c *appengineCache) Set(k string, v interface{}, expire time.Duration) error {
	return c.codec.Set(c.ctx, &memcache.Item{
		Key:        c.key(k),
		Object:     v,
		Expiration: expire,
	})
}

func NewAppEngineCache(request *http.Request) Cache {
	return &appengineCache{
		prefix: os.Getenv("GAE_SERVICE") + "/" + os.Getenv("GAE_VERSION"),
		codec: memcache.Codec{
			Marshal:   json.Marshal,
			Unmarshal: json.Unmarshal,
		},
		ctx: appengine.NewContext(request),
	}
}
